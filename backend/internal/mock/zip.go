package mock

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractZip extracts a zip file from an io.Reader to the specified destination path
func ExtractZip(reader io.ReaderAt, size int64, destPath string) error {
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		return fmt.Errorf("failed to create zip reader: %w", err)
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Extract files
	for _, file := range zipReader.File {
		err := extractZipFile(file, destPath)
		if err != nil {
			return fmt.Errorf("failed to extract file %s: %w", file.Name, err)
		}
	}

	return nil
}

func extractZipFile(file *zip.File, destPath string) error {
	// Construct target path
	targetPath := filepath.Join(destPath, file.Name)

	// Security check: prevent directory traversal
	if !strings.HasPrefix(targetPath, filepath.Clean(destPath)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", file.Name)
	}

	// Open the file in the zip archive
	rc, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file in archive: %w", err)
	}
	defer rc.Close()

	// Check if it's a directory
	if file.FileInfo().IsDir() {
		return os.MkdirAll(targetPath, file.FileInfo().Mode())
	}

	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directories: %w", err)
	}

	// Create the target file
	outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer outFile.Close()

	// Copy file contents
	_, err = io.Copy(outFile, rc)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}
