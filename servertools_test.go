package servertools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ── FileExists ────────────────────────────────────────────────────────────────

func TestFileExists_ExistingFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "testfile-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	f.Close()

	if !FileExists(f.Name()) {
		t.Errorf("expected FileExists(%q) = true, got false", f.Name())
	}
}

func TestFileExists_NonExistentFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.txt")
	if FileExists(path) {
		t.Errorf("expected FileExists(%q) = false, got true", path)
	}
}

func TestFileExists_Directory(t *testing.T) {
	dir := t.TempDir()
	if FileExists(dir) {
		t.Errorf("expected FileExists(%q) = false for a directory, got true", dir)
	}
}

func TestFileExists_EmptyPath(t *testing.T) {
	if FileExists("") {
		t.Error("expected FileExists(\"\") = false, got true")
	}
}

// ── ExpandTilde ───────────────────────────────────────────────────────────────

func TestExpandTilde_WithTilde(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home directory:", err)
	}

	result := ExpandTilde("~/documents/file.txt")
	expected := filepath.Join(home, "documents/file.txt")

	// filepath.Join cleans the path; strings.Replace does not — normalise both.
	if filepath.Clean(result) != filepath.Clean(expected) {
		t.Errorf("ExpandTilde(\"~/documents/file.txt\") = %q, want %q", result, expected)
	}
}

func TestExpandTilde_WithoutTilde(t *testing.T) {
	path := "/absolute/path/to/file.txt"
	result := ExpandTilde(path)
	if result != path {
		t.Errorf("ExpandTilde(%q) = %q, want %q (unchanged)", path, result, path)
	}
}

func TestExpandTilde_TildeOnly(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home directory:", err)
	}

	result := ExpandTilde("~")
	if result != home {
		t.Errorf("ExpandTilde(\"~\") = %q, want %q", result, home)
	}
}

func TestExpandTilde_RelativePath(t *testing.T) {
	path := "relative/path/file.txt"
	result := ExpandTilde(path)
	if result != path {
		t.Errorf("ExpandTilde(%q) = %q, want %q (unchanged)", path, result, path)
	}
}

func TestExpandTilde_TildeInMiddle(t *testing.T) {
	path := "/some/~/path"
	result := ExpandTilde(path)
	if result != path {
		t.Errorf("ExpandTilde(%q) = %q, want %q (tilde mid-path should not expand)", path, result, path)
	}
}

func TestExpandTilde_EmptyString(t *testing.T) {
	result := ExpandTilde("")
	if result != "" {
		t.Errorf("ExpandTilde(\"\") = %q, want \"\"", result)
	}
}

// ── GetFileContentType ────────────────────────────────────────────────────────

func TestGetFileContentType_PlainText(t *testing.T) {
	// http.DetectContentType reads up to 512 bytes; give it a full buffer of
	// printable ASCII so it confidently returns "text/plain" rather than
	// falling back to "application/octet-stream".
	content := []byte(strings.Repeat("Hello, world! This is plain text content. ", 13))
	f := writeTempFile(t, content)
	defer f.Close()

	ct, err := GetFileContentType(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("GetFileContentType() = %q, want prefix \"text/plain\"", ct)
	}
}

func TestGetFileContentType_PNG(t *testing.T) {
	// Minimal valid PNG header (first 8 bytes are the PNG signature).
	pngHeader := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, // IHDR length
		0x49, 0x48, 0x44, 0x52, // "IHDR"
	}
	f := writeTempFile(t, pngHeader)
	defer f.Close()

	ct, err := GetFileContentType(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ct != "image/png" {
		t.Errorf("GetFileContentType() = %q, want \"image/png\"", ct)
	}
}

func TestGetFileContentType_PDF(t *testing.T) {
	pdfHeader := []byte("%PDF-1.4 test content padding padding padding padding")
	f := writeTempFile(t, pdfHeader)
	defer f.Close()

	ct, err := GetFileContentType(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ct != "application/pdf" {
		t.Errorf("GetFileContentType() = %q, want \"application/pdf\"", ct)
	}
}

func TestGetFileContentType_SeeksBackToStart(t *testing.T) {
	content := []byte(strings.Repeat("Hello, world! This is plain text content. ", 13))
	f := writeTempFile(t, content)
	defer f.Close()

	// Call twice — second call must succeed, proving the seek-back worked.
	if _, err := GetFileContentType(f); err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	if _, err := GetFileContentType(f); err != nil {
		t.Errorf("second call failed (seek-back broken): %v", err)
	}
}

func TestGetFileContentType_EmptyFile(t *testing.T) {
	f := writeTempFile(t, []byte{})
	defer f.Close()

	_, err := GetFileContentType(f)
	if err == nil {
		t.Error("expected an error for empty file, got nil")
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

// writeTempFile creates a named temp file with the given content and returns an
// open *os.File positioned at the start, ready for reading.
func writeTempFile(t *testing.T, content []byte) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "content-type-test-*")
	if err != nil {
		t.Fatalf("writeTempFile: create: %v", err)
	}
	if len(content) > 0 {
		if _, err := f.Write(content); err != nil {
			t.Fatalf("writeTempFile: write: %v", err)
		}
		if _, err := f.Seek(0, 0); err != nil {
			t.Fatalf("writeTempFile: seek: %v", err)
		}
	}
	return f
}
