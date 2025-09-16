package mock

import (
	"archive/zip"
	"bytes"
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
		return os.MkdirAll(targetPath, 0755)
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

// CreateZipFromDirectory creates a zip file from a directory and returns it as an io.Reader
// The zip file name will be based on the directory name
func CreateZipFromDirectory(dirPath string) (io.Reader, string, error) {
	// Get the directory name for the zip file name
	dirName := filepath.Base(dirPath)
	zipName := dirName + ".zip"

	// Create a buffer to write the zip file to
	var buf bytes.Buffer

	// Create a new zip writer
	zipWriter := zip.NewWriter(&buf)

	// Walk through the directory
	err := filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if filePath == dirPath {
			return nil
		}

		// Get the relative path from the source directory
		relPath, err := filepath.Rel(dirPath, filePath)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Convert Windows path separators to forward slashes for zip compatibility
		relPath = strings.ReplaceAll(relPath, "\\", "/")

		if info.IsDir() {
			// Add directory entry (with trailing slash)
			_, err := zipWriter.Create(relPath + "/")
			if err != nil {
				return fmt.Errorf("failed to create directory entry: %w", err)
			}
		} else {
			// Add file entry
			err := addFileToZip(zipWriter, filePath, relPath)
			if err != nil {
				return fmt.Errorf("failed to add file %s: %w", filePath, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, "", fmt.Errorf("failed to walk directory: %w", err)
	}

	// Close the zip writer to finalize the archive
	err = zipWriter.Close()
	if err != nil {
		return nil, "", fmt.Errorf("failed to close zip writer: %w", err)
	}

	// Return a reader for the buffer
	return bytes.NewReader(buf.Bytes()), zipName, nil
}

func addFileToZip(zipWriter *zip.Writer, filePath, zipPath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info for metadata
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Create zip file header
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("failed to create file header: %w", err)
	}

	// Set the name in the zip file
	header.Name = zipPath

	// Set compression method
	header.Method = zip.Deflate

	// Create the file in the zip archive
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("failed to create zip entry: %w", err)
	}

	// Copy file contents to zip
	_, err = io.Copy(writer, file)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}
