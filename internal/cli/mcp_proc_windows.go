//go:build windows

package cli

import (
	"os/exec"
)

// setProcAttr sets Windows-specific process attributes
func setProcAttr(cmd *exec.Cmd) {
	// Windows doesn't use Setpgid, so this is a no-op
}