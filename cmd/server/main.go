// @title marchive API
// @version 1.0.0
// @description A secure, fast, and feature-rich Go-Chi backend with JWT authentication, email verification, and password management. Built with modern Go practices and comprehensive security features.
// @termsOfService https://github.com/akramboussanni/marchive/blob/main/LICENSE

// @contact.name API Support
// @contact.url https://github.com/akramboussanni/marchive/issues
// @contact.email support@example.com

// @license.name MIT License
// @license.url https://github.com/akramboussanni/marchive/blob/main/LICENSE

// @host localhost:9520
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name session
// @description JWT session cookie for authenticated endpoints. Automatically set by login endpoint. Required for endpoints marked with @Security CookieAuth.

// @securityDefinitions.apikey RecaptchaToken
// @in header
// @name X-Recaptcha-Token
// @description reCAPTCHA verification token for bot protection. Optional - only required if reCAPTCHA is configured in environment variables. Obtain from reCAPTCHA widget.

// @tag.name Authentication
// @tag.description User registration, login, and token management endpoints. reCAPTCHA verification is optional if configured.

// @tag.name Account
// @tag.description User profile and account management endpoints. All endpoints require session cookie authentication.

// @tag.name Email Verification
// @tag.description Email confirmation and verification endpoints. reCAPTCHA verification is optional if configured.

// @tag.name Password Management
// @tag.description Password reset, change, and recovery endpoints. Public endpoints have optional reCAPTCHA, authenticated endpoints require session cookie.

package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/akramboussanni/marchive/config"
	"github.com/akramboussanni/marchive/internal/api/routes"
	"github.com/akramboussanni/marchive/internal/db"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/repo"
	"github.com/akramboussanni/marchive/internal/services"
	"github.com/akramboussanni/marchive/internal/utils"
)

func createDefaultAdmin(repos *repo.Repos) {
	adminUser, err := repos.User.GetUserByUsername(context.Background(), "admin")
	if err == nil && adminUser != nil {
		log.Println("Admin user already exists, skipping creation")
		return
	}

	passwordBytes := make([]byte, 16)
	_, err = rand.Read(passwordBytes)
	if err != nil {
		log.Printf("Failed to generate random password: %v", err)
		return
	}
	password := base64.URLEncoding.EncodeToString(passwordBytes)

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return
	}

	adminUser = &model.User{
		ID:           utils.GenerateSnowflakeID(),
		Username:     "admin",
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC().Unix(),
		Role:         "admin",
	}

	err = repos.User.CreateUser(context.Background(), adminUser)
	if err != nil {
		log.Printf("Failed to create admin user: %v", err)
		return
	}

	log.Printf("Default admin user created successfully!")
	log.Printf("Username: admin")
	log.Printf("Password: %s", password)
	log.Printf("IMPORTANT: Change this password after first login!")
}

func main() {
	config.Init()

	err := utils.InitSnowflake(1)
	if err != nil {
		panic(err)
	}

	db.Init(config.App.DbConnectionString)
	db.RunMigrations()

	repos := repo.NewRepos(db.DB)

	createDefaultAdmin(repos)

	downloadService := services.NewDownloadService(repos, "/app/downloads", config.App.AnnasApiKey)
	ctx, cancel := context.WithCancel(context.Background())
	go downloadService.StartService(ctx)

	r := routes.SetupRouter(repos)

	port := strconv.Itoa(config.App.AppPort)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	if config.App.TLSEnabled {
		if config.App.TLSCertFile == "" || config.App.TLSKeyFile == "" {
			log.Fatal("TLS_ENABLED is true but TLS_CERT_FILE or TLS_KEY_FILE is not set")
		}

		server.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		}
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("shutting down server...")
		cancel()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("server forced to shutdown: %v", err)
		}
		log.Println("server exited gracefully")
	}()

	protocol := "http"
	if config.App.TLSEnabled {
		protocol = "https"
	}

	log.Printf("server running @ %s://localhost:%s", protocol, port)

	if config.App.TLSEnabled {
		if err := server.ListenAndServeTLS(config.App.TLSCertFile, config.App.TLSKeyFile); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error when starting TLS server: %v", err)
		}
	} else {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error when starting server: %v", err)
		}
	}
}
