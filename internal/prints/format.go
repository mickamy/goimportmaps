package prints

import (
	"fmt"
)

type Format string

const (
	FormatGraphviz Format = "graphviz"
	FormatHTML     Format = "html"
	FormatMermaid  Format = "mermaid"
	FormatText     Format = "text"
)

func NewFormat(format string) (Format, error) {
	switch format {
	case string(FormatGraphviz):
		return FormatGraphviz, nil
	case string(FormatHTML):
		return FormatHTML, nil
	case string(FormatMermaid):
		return FormatMermaid, nil
	case string(FormatText):
		return FormatText, nil
	default:
		return "", fmt.Errorf("unsupported format %s", format)
	}
}
