package export

import (
	"context"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConvertHTMLToPDF_Success(t *testing.T) {
	expectedPDF := []byte("%PDF-1.4 fake pdf content")
	inputHTML := "<html><body><h1>Hello</h1></body></html>"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/forms/chromium/convert/html" {
			t.Errorf("expected /forms/chromium/convert/html, got %s", r.URL.Path)
		}

		// Verify multipart form
		contentType := r.Header.Get("Content-Type")
		_, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			t.Fatalf("parsing content-type: %v", err)
		}

		reader := multipart.NewReader(r.Body, params["boundary"])
		part, err := reader.NextPart()
		if err != nil {
			t.Fatalf("reading part: %v", err)
		}

		if part.FormName() != "files" {
			t.Errorf("expected field name 'files', got %q", part.FormName())
		}
		if part.FileName() != "index.html" {
			t.Errorf("expected filename 'index.html', got %q", part.FileName())
		}
		if ct := part.Header.Get("Content-Type"); ct != "text/html" {
			t.Errorf("expected content-type 'text/html', got %q", ct)
		}

		body, _ := io.ReadAll(part)
		if string(body) != inputHTML {
			t.Errorf("expected HTML %q, got %q", inputHTML, string(body))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(expectedPDF)
	}))
	defer server.Close()

	client := NewGotenbergClient(server.URL)
	pdf, err := client.ConvertHTMLToPDF(context.Background(), inputHTML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(pdf) != string(expectedPDF) {
		t.Errorf("expected %q, got %q", expectedPDF, pdf)
	}
}

func TestConvertHTMLToPDF_GotenbergError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("invalid file format"))
	}))
	defer server.Close()

	client := NewGotenbergClient(server.URL)
	_, err := client.ConvertHTMLToPDF(context.Background(), "<html></html>")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := "gotenberg returned status 400: invalid file format"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestConvertHTMLToPDF_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Block until context is cancelled — the client should abort
		<-r.Context().Done()
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := NewGotenbergClient(server.URL)
	_, err := client.ConvertHTMLToPDF(ctx, "<html></html>")
	if err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}
}

func TestConvertHTMLToPDF_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("chromium crashed"))
	}))
	defer server.Close()

	client := NewGotenbergClient(server.URL)
	_, err := client.ConvertHTMLToPDF(context.Background(), "<html></html>")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := "gotenberg returned status 500: chromium crashed"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

