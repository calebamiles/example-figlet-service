package figlet

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	defaultBrewFigletLocation   = "/usr/local/bin/figlet"
	defaultUbuntuFigletLocation = "/usr/bin/figlet"
	defaultFedoraFigletLocation = "/usr/bin/figlet"
)

// A Transformer applies a Figlet transformation to a given string
type Transformer interface {
	Figletize(string) (string, error)
}

// NewTransformer returns a new Figlet transformer
func NewTransformer() Transformer {
	return &execFigletTransformer{}
}

type execFigletTransformer struct{}

func (t *execFigletTransformer) Figletize(in string) (string, error) {
	var cmd *exec.Cmd

	if _, existErr := os.Stat(defaultUbuntuFigletLocation); existErr == nil {
		cmd = exec.Command(defaultUbuntuFigletLocation)
	} else if _, existErr := os.Stat(defaultFedoraFigletLocation); existErr == nil {
		cmd = exec.Command(defaultFedoraFigletLocation)
	} else if _, existErr := os.Stat(defaultBrewFigletLocation); existErr == nil {
		cmd = exec.Command(defaultBrewFigletLocation)
	} else if _, pathErr := exec.LookPath("figlet"); pathErr == nil {
		cmd = exec.Command("figlet")
	}

	if cmd == nil {
		return "", fmt.Errorf("unable to locate figlet executable")
	}

	cmd.Stdin = strings.NewReader(in)

	figletedOutput, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(figletedOutput), nil
}
