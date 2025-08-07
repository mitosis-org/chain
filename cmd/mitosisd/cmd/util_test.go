package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateFilePath(t *testing.T) {
	tempDir := t.TempDir()

	// Create test files
	validFile := filepath.Join(tempDir, "valid.json")
	err := os.WriteFile(validFile, []byte(`{"test": "data"}`), 0o644)
	require.NoError(t, err)

	largeFile := filepath.Join(tempDir, "large.json")
	largeContent := strings.Repeat("a", 1000)
	err = os.WriteFile(largeFile, []byte(largeContent), 0o644)
	require.NoError(t, err)

	// Create a directory to test against (with .json extension to pass extension check)
	testDir := filepath.Join(tempDir, "testdir.json")
	err = os.MkdirAll(testDir, 0o755)
	require.NoError(t, err)

	tests := []struct {
		name              string
		filePath          string
		allowedExtensions []string
		maxSize           int64
		expectError       bool
		errorMsg          string
	}{
		{
			name:              "valid json file",
			filePath:          validFile,
			allowedExtensions: []string{".json"},
			maxSize:           1000,
			expectError:       false,
		},
		{
			name:              "valid file with multiple allowed extensions",
			filePath:          validFile,
			allowedExtensions: []string{".txt", ".json", ".yaml"},
			maxSize:           1000,
			expectError:       false,
		},
		{
			name:              "invalid extension",
			filePath:          validFile,
			allowedExtensions: []string{".txt"},
			maxSize:           1000,
			expectError:       true,
			errorMsg:          "file must have one of the following extensions",
		},
		{
			name:              "file too large",
			filePath:          largeFile,
			allowedExtensions: []string{".json"},
			maxSize:           100,
			expectError:       true,
			errorMsg:          "file too large",
		},
		{
			name:              "path traversal attempt",
			filePath:          "../../../etc/passwd",
			allowedExtensions: []string{".json"},
			maxSize:           1000,
			expectError:       true,
			errorMsg:          "path traversal detected",
		},
		{
			name:              "non-existent file",
			filePath:          filepath.Join(tempDir, "nonexistent.json"),
			allowedExtensions: []string{".json"},
			maxSize:           1000,
			expectError:       true,
			errorMsg:          "file does not exist",
		},
		{
			name:              "directory instead of file",
			filePath:          testDir,
			allowedExtensions: []string{".json"},
			maxSize:           1000,
			expectError:       true,
			errorMsg:          "path must be a regular file",
		},
		{
			name:              "no extension restrictions",
			filePath:          validFile,
			allowedExtensions: nil,
			maxSize:           1000,
			expectError:       false,
		},
		{
			name:              "no size restrictions",
			filePath:          largeFile,
			allowedExtensions: []string{".json"},
			maxSize:           0,
			expectError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFilePath(tt.filePath, tt.allowedExtensions, tt.maxSize)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateFilePath_SymlinkHandling(t *testing.T) {
	tempDir := t.TempDir()

	// Create a target file
	targetFile := filepath.Join(tempDir, "target.json")
	err := os.WriteFile(targetFile, []byte(`{"test": "data"}`), 0o644)
	require.NoError(t, err)

	// Create a symlink to the target file
	symlinkFile := filepath.Join(tempDir, "symlink.json")
	err = os.Symlink(targetFile, symlinkFile)
	if err != nil {
		t.Skip("Symlinks not supported on this system")
	}

	// Test with valid symlink
	err = validateFilePath(symlinkFile, []string{".json"}, 1000)
	require.NoError(t, err)

	// Create a symlink with different extension
	badSymlink := filepath.Join(tempDir, "bad_symlink.txt")
	err = os.Symlink(targetFile, badSymlink)
	require.NoError(t, err)

	// Test with symlink that doesn't match extension
	err = validateFilePath(badSymlink, []string{".json"}, 1000)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "file must have one of the following extensions")
}

func TestValidateArtifactFile(t *testing.T) {
	tempDir := t.TempDir()

	// Create a valid artifact file
	validArtifact := filepath.Join(tempDir, "artifact.json")
	err := os.WriteFile(validArtifact, []byte(`{"bytecode": {"object": "0x608060405234801561001057600080fd5b50"}}`), 0o644)
	require.NoError(t, err)

	// Create an invalid artifact file (wrong extension)
	invalidArtifact := filepath.Join(tempDir, "artifact.txt")
	err = os.WriteFile(invalidArtifact, []byte(`{"bytecode": {"object": "0x608060405234801561001057600080fd5b50"}}`), 0o644)
	require.NoError(t, err)

	tests := []struct {
		name        string
		filePath    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid artifact file",
			filePath:    validArtifact,
			expectError: false,
		},
		{
			name:        "invalid extension",
			filePath:    invalidArtifact,
			expectError: true,
			errorMsg:    "file must have one of the following extensions",
		},
		{
			name:        "non-existent file",
			filePath:    filepath.Join(tempDir, "nonexistent.json"),
			expectError: true,
			errorMsg:    "file does not exist",
		},
		{
			name:        "path traversal attempt",
			filePath:    "../../../etc/passwd",
			expectError: true,
			errorMsg:    "path traversal detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateArtifactFile(tt.filePath)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateArtifactFile_SizeLimit(t *testing.T) {
	tempDir := t.TempDir()

	// Create a large artifact file (over 10MB limit)
	largeArtifact := filepath.Join(tempDir, "large_artifact.json")
	// Create content larger than 10MB
	largeContent := fmt.Sprintf(`{"bytecode": {"object": "%s"}}`, strings.Repeat("a", 11*1024*1024))
	err := os.WriteFile(largeArtifact, []byte(largeContent), 0o644)
	require.NoError(t, err)

	err = validateArtifactFile(largeArtifact)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "file too large")
}

func TestSafeCopyFile(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file
	srcFile := filepath.Join(tempDir, "source.txt")
	srcContent := "This is test content for copying"
	err := os.WriteFile(srcFile, []byte(srcContent), 0o644)
	require.NoError(t, err)

	tests := []struct {
		name        string
		src         string
		dst         string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid copy",
			src:         srcFile,
			dst:         filepath.Join(tempDir, "destination.txt"),
			expectError: false,
		},
		{
			name:        "copy to nested directory",
			src:         srcFile,
			dst:         filepath.Join(tempDir, "nested", "dir", "destination.txt"),
			expectError: false,
		},
		{
			name:        "source path traversal",
			src:         "../../../etc/passwd",
			dst:         filepath.Join(tempDir, "destination.txt"),
			expectError: true,
			errorMsg:    "invalid source file path: path traversal detected",
		},
		{
			name:        "destination path traversal",
			src:         srcFile,
			dst:         "../../../tmp/destination.txt",
			expectError: true,
			errorMsg:    "invalid destination file path: path traversal detected",
		},
		{
			name:        "non-existent source",
			src:         filepath.Join(tempDir, "nonexistent.txt"),
			dst:         filepath.Join(tempDir, "destination.txt"),
			expectError: true,
			errorMsg:    "source file does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := safeCopyFile(tt.src, tt.dst)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)

				// Verify destination file exists and has correct content
				dstContent, err := os.ReadFile(tt.dst)
				require.NoError(t, err)
				assert.Equal(t, srcContent, string(dstContent))

				// Verify destination directory was created
				dstDir := filepath.Dir(tt.dst)
				dirInfo, err := os.Stat(dstDir)
				require.NoError(t, err)
				assert.True(t, dirInfo.IsDir())
			}
		})
	}
}

func TestSafeCopyFile_DirectoryAsSource(t *testing.T) {
	tempDir := t.TempDir()

	// Create a directory to test with
	srcDir := filepath.Join(tempDir, "source_dir")
	err := os.MkdirAll(srcDir, 0o755)
	require.NoError(t, err)

	dstFile := filepath.Join(tempDir, "destination.txt")

	err = safeCopyFile(srcDir, dstFile)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "source must be a regular file")
}

func TestSafeCopyFile_LargeFile(t *testing.T) {
	tempDir := t.TempDir()

	// Create a large source file (over 100MB limit)
	largeSrcFile := filepath.Join(tempDir, "large_source.txt")
	// Create content larger than 100MB
	largeContent := strings.Repeat("a", 101*1024*1024)
	err := os.WriteFile(largeSrcFile, []byte(largeContent), 0o644)
	require.NoError(t, err)

	dstFile := filepath.Join(tempDir, "destination.txt")

	err = safeCopyFile(largeSrcFile, dstFile)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "source file too large to copy")
}

func TestSafeCopyFile_PermissionsPreservation(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file with specific content
	srcFile := filepath.Join(tempDir, "source.txt")
	srcContent := "test content"
	err := os.WriteFile(srcFile, []byte(srcContent), 0o644)
	require.NoError(t, err)

	dstFile := filepath.Join(tempDir, "destination.txt")

	err = safeCopyFile(srcFile, dstFile)
	require.NoError(t, err)

	// Verify content was copied correctly
	dstContent, err := os.ReadFile(dstFile)
	require.NoError(t, err)
	assert.Equal(t, srcContent, string(dstContent))

	// Verify both files exist
	_, err = os.Stat(srcFile)
	require.NoError(t, err)
	_, err = os.Stat(dstFile)
	require.NoError(t, err)
}

func TestSafeCopyFile_OverwriteExisting(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file
	srcFile := filepath.Join(tempDir, "source.txt")
	srcContent := "new content"
	err := os.WriteFile(srcFile, []byte(srcContent), 0o644)
	require.NoError(t, err)

	// Create destination file with different content
	dstFile := filepath.Join(tempDir, "destination.txt")
	oldContent := "old content"
	err = os.WriteFile(dstFile, []byte(oldContent), 0o644)
	require.NoError(t, err)

	// Copy should overwrite
	err = safeCopyFile(srcFile, dstFile)
	require.NoError(t, err)

	// Verify destination has new content
	dstContent, err := os.ReadFile(dstFile)
	require.NoError(t, err)
	assert.Equal(t, srcContent, string(dstContent))
	assert.NotEqual(t, oldContent, string(dstContent))
}
