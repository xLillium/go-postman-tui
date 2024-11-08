package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/xlillium/go-postman-tui/formatters"
)

// PerformRequest performs an HTTP request with the specified method, URL, and body.
func PerformRequest(method, url, body string) (formattedResponse string, status string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if url == "" {
		return "", "", fmt.Errorf("URL cannot be empty")
	}

	var req *http.Request

	// Prepare request based on method
	switch method {
	case "GET":
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	case "POST":
		if !isValidJSON(body) {
			return "", "", fmt.Errorf("Invalid JSON in request body")
		}
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	default:
		return "", "", fmt.Errorf("Unsupported method: %s", method)
	}

	if err != nil {
		return "", "", fmt.Errorf("Error creating request: %v", err)
	}

	// Perform the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("Error performing request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("Error reading response: %v", err)
	}

	// Format the response
	formattedResponse, err = formatResponse(bodyBytes, resp.Header.Get("Content-Type"))
	if err != nil {
		return "", "", fmt.Errorf("Error formatting response: %v", err)
	}

	return formattedResponse, resp.Status, nil
}

// isValidJSON checks if a string is valid JSON.
func isValidJSON(data string) bool {
	var js interface{}
	return json.Unmarshal([]byte(data), &js) == nil
}

// formatResponse formats the response body based on Content-Type.
func formatResponse(body []byte, contentType string) (string, error) {
	var lexer chroma.Lexer

	// Select lexer based on Content-Type
	switch {
	case strings.Contains(contentType, "application/json"):
		lexer = lexers.Get("json")

		// Pretty-print JSON
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, body, "", "  ")
		if err != nil {
			return "", err
		}
		body = prettyJSON.Bytes()

	case strings.Contains(contentType, "application/xml"), strings.Contains(contentType, "text/xml"):
		lexer = lexers.Get("xml")

	case strings.Contains(contentType, "text/html"):
		lexer = lexers.Get("html")

	default:
		lexer = lexers.Fallback
	}

	// Tokenize
	iterator, err := lexer.Tokenise(nil, string(body))
	if err != nil {
		return "", err
	}

	// Format using the custom formatter
	formatter := &formatters.TviewFormatter{}
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	var buff bytes.Buffer
	err = formatter.Format(&buff, style, iterator)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}

