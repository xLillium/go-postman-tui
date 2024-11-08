package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"

	"github.com/xlillium/go-postman-tui/formatters"
)

// PerformGetRequest performs a GET request with a timeout and updates the UI
func PerformGetRequest(url string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if url == "" {
		url = "http://api.github.com"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", fmt.Errorf("Error creating request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("Error: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("Error reading response: %v", err)
	}

	// Format the response
	formattedResponse, err := formatResponse(bodyBytes, resp.Header.Get("Content-Type"))
	if err != nil {
		return "", "", fmt.Errorf("Error formatting response: %v", err)
	}

	return formattedResponse, resp.Status, nil
}

// Similar changes for PerformRequest
func PerformRequest(method, url, body string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if url == "" {
		url = "http://api.github.com"
	}

	var req *http.Request
	var err error

	if method == "GET" {
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	} else if method == "POST" {
		var jsonBody interface{}
		err := json.Unmarshal([]byte(body), &jsonBody)
		if err != nil {
			return "", "", fmt.Errorf("Invalid JSON: %v", err)
		}
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		return "", "", fmt.Errorf("Unsupported method: %s", method)
	}

	if err != nil {
		return "", "", fmt.Errorf("Error creating request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("Error: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("Error reading response: %v", err)
	}

	// Format the response
	formattedResponse, err := formatResponse(bodyBytes, resp.Header.Get("Content-Type"))
	if err != nil {
		return "", "", fmt.Errorf("Error formatting response: %v", err)
	}

	return formattedResponse, resp.Status, nil
}

func formatResponse(body []byte, contentType string) (string, error) {
	var lexer chroma.Lexer
	var err error

	// Select lexer based on Content-Type
	switch {
	case strings.Contains(contentType, "application/json"):
		lexer = lexers.Get("json")

		// Pretty-print JSON
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, body, "", "  ")
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
