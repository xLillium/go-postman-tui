package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"bytes"
	"encoding/json"
	"github.com/rivo/tview"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"

	"github.com/xlillium/go-postman-tui/formatters"
)

// PerformGetRequest performs a GET request with a timeout and updates the UI
func PerformGetRequest(url string, middleBox, bottomBox *tview.TextView, app *tview.Application) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		app.QueueUpdateDraw(func() {
			bottomBox.SetText(fmt.Sprintf("[red]Error creating request: %v", err))
			middleBox.SetText("")
		})
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		app.QueueUpdateDraw(func() {
			bottomBox.SetText(fmt.Sprintf("[red]Error: %v", err))
			middleBox.SetText("")
		})
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		app.QueueUpdateDraw(func() {
			bottomBox.SetText(fmt.Sprintf("[red]Error reading response: %v", err))
			middleBox.SetText("")
		})
		return
	}

	// bodyString := string(bodyBytes)

	// Format the response
	formattedResponse, err := formatResponse(bodyBytes, resp.Header.Get("Content-Type"))
	if err != nil {
		app.QueueUpdateDraw(func() {
			bottomBox.SetText(fmt.Sprintf("[red]Error formatting response: %v", err))
			middleBox.SetText("")
		})
		return
	}

	app.QueueUpdateDraw(func() {
		middleBox.SetText(formattedResponse)
		bottomBox.SetText(fmt.Sprintf("[green]Request successful (%s)", resp.Status))
	})

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
