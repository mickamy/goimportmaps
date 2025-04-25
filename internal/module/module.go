package module

import (
	"bytes"
	"os/exec"
	"strings"
)

// Path returns the current module name (e.g., github.com/xxx/yyy).
func Path() (string, error) {
	cmd := exec.Command("go", "list", "-m")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
