package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// validateFilePath validates and sanitizes a file path to prevent path traversal attacks
func validateFilePath(filePath string, allowedExtensions []string, maxSize int64) error {
	// Clean the path to resolve any ".." elements
	cleanPath := filepath.Clean(filePath)

	// Check if the path contains any ".." after cleaning (indicates path traversal attempt)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected in file path")
	}

	// Validate file extension if specified
	if len(allowedExtensions) > 0 {
		hasValidExt := false
		lowerPath := strings.ToLower(cleanPath)
		for _, ext := range allowedExtensions {
			if strings.HasSuffix(lowerPath, strings.ToLower(ext)) {
				hasValidExt = true
				break
			}
		}
		if !hasValidExt {
			return fmt.Errorf("file must have one of the following extensions: %v", allowedExtensions)
		}
	}

	// Resolve any symbolic links to get the actual file path
	realPath, err := filepath.EvalSymlinks(cleanPath)
	if err != nil {
		// If EvalSymlinks fails, check if file exists normally
		if _, statErr := os.Stat(cleanPath); os.IsNotExist(statErr) {
			return fmt.Errorf("file does not exist: %s", cleanPath)
		}
		// Use the original clean path if EvalSymlinks fails but file exists
		realPath = cleanPath
	}

	// Re-validate file extension after resolving symlinks (prevents symlink attacks)
	if len(allowedExtensions) > 0 {
		hasValidExt := false
		lowerRealPath := strings.ToLower(realPath)
		for _, ext := range allowedExtensions {
			if strings.HasSuffix(lowerRealPath, strings.ToLower(ext)) {
				hasValidExt = true
				break
			}
		}
		if !hasValidExt {
			return fmt.Errorf("file (resolved from symlink) must have one of the following extensions: %v", allowedExtensions)
		}
	}

	// Check if file exists and is readable
	fileInfo, err := os.Stat(realPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", cleanPath)
	}
	if err != nil {
		return fmt.Errorf("cannot access file: %w", err)
	}

	// Ensure it's a regular file (not a directory or special file)
	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("path must be a regular file, not a directory or special file")
	}

	// Check file size if maxSize is specified (prevent extremely large files that could cause DoS)
	if maxSize > 0 && fileInfo.Size() > maxSize {
		return fmt.Errorf("file too large (max %d bytes)", maxSize)
	}

	return nil
}

// validateArtifactFile validates and sanitizes the artifact file path to prevent path traversal attacks
func validateArtifactFile(artifactFile string) error {
	const maxFileSize = 10 * 1024 * 1024 // 10MB limit
	allowedExtensions := []string{".json"}
	return validateFilePath(artifactFile, allowedExtensions, maxFileSize)
}

// safeCopyFile safely copies a file from src to dst with path validation
func safeCopyFile(src, dst string) error {
	// Validate and clean source path
	cleanSrc := filepath.Clean(src)
	if strings.Contains(cleanSrc, "..") {
		return fmt.Errorf("invalid source file path: path traversal detected")
	}

	// Validate and clean destination path
	cleanDst := filepath.Clean(dst)
	if strings.Contains(cleanDst, "..") {
		return fmt.Errorf("invalid destination file path: path traversal detected")
	}

	// Check if source file exists and is accessible
	srcInfo, err := os.Stat(cleanSrc)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("source file does not exist: %s", cleanSrc)
		}
		return fmt.Errorf("cannot access source file: %w", err)
	}

	// Ensure source is a regular file
	if !srcInfo.Mode().IsRegular() {
		return fmt.Errorf("source must be a regular file, not a directory or special file")
	}

	// Check file size (prevent copying extremely large files)
	const maxCopySize = 100 * 1024 * 1024 // 100MB limit
	if srcInfo.Size() > maxCopySize {
		return fmt.Errorf("source file too large to copy (max %d bytes)", maxCopySize)
	}

	// Ensure destination directory exists
	dstDir := filepath.Dir(cleanDst)
	if err := os.MkdirAll(dstDir, 0o700); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	sourceFile, err := os.Open(cleanSrc)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(cleanDst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}
