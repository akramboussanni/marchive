package utils

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	lowerRegex = regexp.MustCompile(`[a-z]`)
	upperRegex = regexp.MustCompile(`[A-Z]`)
	digitRegex = regexp.MustCompile(`\d`)
	specialRegex = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
)

// PasswordRequirements defines the password validation rules
type PasswordRequirements struct {
	MinLength     int
	RequireLower  bool
	RequireUpper  bool
	RequireDigit  bool
	RequireSpecial bool
}

// DefaultPasswordRequirements returns the standard password requirements
func DefaultPasswordRequirements() PasswordRequirements {
	return PasswordRequirements{
		MinLength:     8,
		RequireLower:  true,
		RequireUpper:  true,
		RequireDigit:  true,
		RequireSpecial: false,
	}
}

// ValidatePasswordWithDetails validates a password and returns detailed error information
func ValidatePasswordWithDetails(password string, req PasswordRequirements) (bool, []string) {
	var errors []string
	
	if len(password) < req.MinLength {
		errors = append(errors, fmt.Sprintf("Password must be at least %d characters long", req.MinLength))
	}
	
	if req.RequireLower && !lowerRegex.MatchString(password) {
		errors = append(errors, "Password must contain at least one lowercase letter")
	}
	
	if req.RequireUpper && !upperRegex.MatchString(password) {
		errors = append(errors, "Password must contain at least one uppercase letter")
	}
	
	if req.RequireDigit && !digitRegex.MatchString(password) {
		errors = append(errors, "Password must contain at least one number")
	}
	
	if req.RequireSpecial && !specialRegex.MatchString(password) {
		errors = append(errors, "Password must contain at least one special character")
	}
	
	return len(errors) == 0, errors
}

// GetPasswordRequirementsText returns a human-readable description of password requirements
func GetPasswordRequirementsText(req PasswordRequirements) string {
	var requirements []string
	
	requirements = append(requirements, fmt.Sprintf("At least %d characters", req.MinLength))
	
	if req.RequireLower {
		requirements = append(requirements, "One lowercase letter")
	}
	
	if req.RequireUpper {
		requirements = append(requirements, "One uppercase letter")
	}
	
	if req.RequireDigit {
		requirements = append(requirements, "One number")
	}
	
	if req.RequireSpecial {
		requirements = append(requirements, "One special character")
	}
	
	return strings.Join(requirements, ", ")
}

// IsValidPassword validates password using default requirements (backward compatibility)
func IsValidPassword(pw string) bool {
	valid, _ := ValidatePasswordWithDetails(pw, DefaultPasswordRequirements())
	return valid
}
