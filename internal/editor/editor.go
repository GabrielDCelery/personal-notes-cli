package editor

import (
	"os"
	"os/exec"
)

func OpenPathInEditor(notePath string, editor string) error {
	cmd := exec.Command(editor, notePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
