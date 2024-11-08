package formatters

import (
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/chroma"
)

// TviewFormatter is a custom formatter for tview
type TviewFormatter struct{}

// Format formats the tokens using tview's color tags
func (f *TviewFormatter) Format(w io.Writer, style *chroma.Style, iterator chroma.Iterator) error {
	for token := iterator(); token != chroma.EOF; token = iterator() {
		entry := style.Get(token.Type)
		fgColor := entry.Colour.String()

		// Apply tview color tags without quotes
		if fgColor != "" {
			fmt.Fprintf(w, `[%s]`, fgColor)
		}

		// Escape tview special characters
		text := strings.ReplaceAll(token.Value, `[`, `[[`)
		text = strings.ReplaceAll(text, `]`, `]]`)

		fmt.Fprint(w, text)

		// Reset to previous color
		if fgColor != "" {
			fmt.Fprint(w, `[-]`)
		}
	}
	return nil
}

