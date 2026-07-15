package storage

import (
	"io"
	"strings"
	"testing"
)

func TestFilesystemSaveOpenDelete(t *testing.T) {
	store, err := NewFilesystem(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	written, err := store.Save("companies/1/resume.pdf", strings.NewReader("content"))
	if err != nil || written != 7 {
		t.Fatalf("save failed: %d, %v", written, err)
	}
	reader, err := store.Open("companies/1/resume.pdf")
	if err != nil {
		t.Fatal(err)
	}
	data, _ := io.ReadAll(reader)
	_ = reader.Close()
	if string(data) != "content" {
		t.Fatalf("unexpected content %q", data)
	}
	if err := store.Delete("companies/1/resume.pdf"); err != nil {
		t.Fatal(err)
	}
}

func TestFilesystemRejectsTraversal(t *testing.T) {
	store, err := NewFilesystem(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"../escape", "/absolute", "a/../../escape"} {
		if _, err := store.Save(key, strings.NewReader("x")); err == nil {
			t.Fatalf("unsafe key %q was accepted", key)
		}
	}
}
