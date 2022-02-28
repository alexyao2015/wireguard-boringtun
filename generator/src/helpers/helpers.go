package helpers

import (
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func WrapError(f func() ([]byte, error), fatal bool, filename string) ([]byte, error) {
	data, err := f()
	if err != nil {
		log.Error("Error reading %s!", filename)
		if fatal {
			log.WithField("err", err).Error("Cannot continue")
			log.Fatal("Exiting...")
		}
	}
	return data, err
}

func RunCmd(stdin string, prg string, args ...string) (string, error) {
	cmd, sb := exec.Command(prg, args...), new(strings.Builder)
	cmd.Stdout = sb
	cmd.Stdin = strings.NewReader(stdin)
	err := cmd.Run()

	return sb.String(), err
}
