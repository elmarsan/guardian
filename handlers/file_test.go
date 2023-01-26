package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/elmarsan/guardian/files"
)

// newTestStorage returns local storage used for tests
func newTestStorage() (*files.LocalStorage, error) {
	p := "test_files"

	// Create local storage
	storage := files.NewLocalStorage(p)

	// Delete storage dir if exits
	os.RemoveAll(p)

	// Create storage directory
	err := os.Mkdir(p, 0775)
	if err != nil {
		return nil, err
	}

	// Save notes.txt into storage
	f := "notes.txt"
	content := []byte("My notes file")
	b := bytes.NewBuffer(content)

	err = storage.Save(f, b)
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func TestFileDownload(t *testing.T) {
	t.Run("should download file", func(t *testing.T) {
		storage, err := newTestStorage()
		if err != nil {
			t.Fatal(storage)
		}

		handler := NewDownloadFile(storage)
		fname := "test_files/notes.txt"

		req := httptest.NewRequest("GET", "/files/download/"+fname, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		t.Log(rec.Result().Header)

		if rec.Result().StatusCode != http.StatusOK {
			t.Errorf("StatusCode should be 200")
		}

		if rec.Header().Get("Content-Disposition") != "attachment; filename="+fname {
			t.Errorf("Wrong Content-Disposition header")
		}
	})

	t.Run("should return not found when file does not exist", func(t *testing.T) {
		storage, err := newTestStorage()
		if err != nil {
			t.Fatal(storage)
		}

		handler := NewDownloadFile(storage)

		req := httptest.NewRequest("GET", "/files/download/test_files/awesome_stuff.png", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		t.Log(rec.Result().StatusCode)

		if rec.Result().StatusCode != http.StatusNotFound {
			t.Errorf("StatusCode should be 404")
		}
	})
}

func TestFiles(t *testing.T) {
	t.Run("should return status ok", func(t *testing.T) {
		storage, err := newTestStorage()
		if err != nil {
			t.Fatal(storage)
		}

		handler := NewFiles(storage, "../templates/login.tmpl")

		req := httptest.NewRequest("GET", "/files", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusOK {
			t.Errorf("StatusCode should be 200")
		}
	})
}
