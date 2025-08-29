package anna

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func extractMetaInformation(meta string) (language, format, size string) {
	if meta == "" {
		return "", "", ""
	}

	parts := strings.Split(meta, "·")

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)

		if strings.Contains(trimmed, "[") && strings.Contains(trimmed, "]") {
			language = trimmed
		}

		upper := strings.ToUpper(trimmed)
		if upper == "PDF" || upper == "EPUB" || upper == "ZIP" || upper == "MOBI" ||
			upper == "AZW3" || upper == "TXT" || upper == "DOC" || upper == "DOCX" {
			format = trimmed
		}

		if strings.Contains(upper, "MB") || strings.Contains(upper, "KB") || strings.Contains(upper, "GB") {
			size = trimmed
		}
	}

	log.Printf("extractMetaInformation: input='%s' -> language='%s', format='%s', size='%s'",
		meta, language, format, size)

	return language, format, size
}

func (b *Book) Download(secretKey, folderPath string) error {
	apiURL := fmt.Sprintf(AnnasDownloadEndpoint, b.Hash, secretKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var apiResp fastDownloadResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return err
	}
	if apiResp.DownloadURL == "" {
		if apiResp.Error != "" {
			return errors.New(apiResp.Error)
		}
		return errors.New("failed to get download URL")
	}

	downloadResp, err := http.Get(apiResp.DownloadURL)
	if err != nil {
		return err
	}
	defer downloadResp.Body.Close()

	if downloadResp.StatusCode != http.StatusOK {
		return errors.New("failed to download file")
	}

	filename := b.Title + "." + b.Format
	filename = strings.ReplaceAll(filename, "/", "_")
	filePath := filepath.Join(folderPath, filename)

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, downloadResp.Body)
	return err
}

func (b *Book) String() string {
	return fmt.Sprintf("Title: %s\nAuthors: %s\nPublisher: %s\nLanguage: %s\nFormat: %s\nSize: %s\nURL: %s\nHash: %s",
		b.Title, b.Authors, b.Publisher, b.Language, b.Format, b.Size, b.URL, b.Hash)
}

func (b *Book) ToJSON() (string, error) {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func GetBookMetadata(hash string) (*Book, error) {
	log.Printf("Fetching metadata for book hash: %s", hash)

	books, err := FindBook(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to search for book metadata: %w", err)
	}

	for _, book := range books {
		if book.Hash == hash {
			log.Printf("Found metadata for hash %s: %s", hash, book.Title)
			return book, nil
		}
	}

	return nil, fmt.Errorf("book with hash %s not found", hash)
}
