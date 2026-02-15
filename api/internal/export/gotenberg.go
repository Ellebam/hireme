package export

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"
)

// PDFConverter converts HTML to PDF.
type PDFConverter interface {
	ConvertHTMLToPDF(ctx context.Context, html string) ([]byte, error)
}

// GotenbergClient calls Gotenberg's conversion endpoints.
type GotenbergClient struct {
	url        string
	httpClient *http.Client
}

// NewGotenbergClient creates a new Gotenberg client.
func NewGotenbergClient(url string) *GotenbergClient {
	return &GotenbergClient{
		url: url,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ConvertHTMLToPDF sends HTML to Gotenberg and returns PDF bytes.
func (g *GotenbergClient) ConvertHTMLToPDF(ctx context.Context, html string) ([]byte, error) {
	return g.convertHTML(ctx, html, "/forms/chromium/convert/html")
}

// convertHTML sends HTML to a Gotenberg endpoint and returns the converted bytes.
func (g *GotenbergClient) convertHTML(ctx context.Context, html string, endpoint string) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a form file part with explicit text/html content type
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="files"; filename="index.html"`)
	header.Set("Content-Type", "text/html")

	part, err := writer.CreatePart(header)
	if err != nil {
		return nil, fmt.Errorf("creating multipart part: %w", err)
	}
	if _, err := part.Write([]byte(html)); err != nil {
		return nil, fmt.Errorf("writing html to part: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("closing multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.url+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling gotenberg: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gotenberg returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
