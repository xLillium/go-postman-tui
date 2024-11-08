package handlers

import (
    "context"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"

    "github.com/rivo/tview"
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

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        app.QueueUpdateDraw(func() {
            bottomBox.SetText(fmt.Sprintf("[red]Error reading response: %v", err))
            middleBox.SetText("")
        })
        return
    }

    bodyString := string(bodyBytes)

    app.QueueUpdateDraw(func() {
        middleBox.SetText(bodyString)
        bottomBox.SetText(fmt.Sprintf("[green]Request successful (%s)", resp.Status))
    })
}

