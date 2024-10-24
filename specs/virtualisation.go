package specs

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/2zqa/ssot-specs-collector/aggregator"
)

const vmDetectorBinary = "systemd-detect-virt"

type Virtualization struct {
	// Type of virtualization (if any), e.g. physical/kvm/lxc
	Type string `json:"type"`
}

func FetchVirtualization() (Virtualization, error) {
	a := aggregator.New("virtualization")
	cmd := exec.Command(vmDetectorBinary)

	output, err := cmd.Output()
	if err != nil {
		exitError, isExitError := err.(*exec.ExitError)
		if a.CheckPredicateAndAdd(!isExitError, err.Error()) {
			return Virtualization{}, a.ErrorOrNil()
		}

		exitCode := exitError.ExitCode()
		if a.CheckPredicateAndAdd(exitCode != 0 && exitCode != 1, fmt.Sprintf("unexpected error code %d", exitCode)) {
			// Unknown exit code, give up
			return Virtualization{}, a.ErrorOrNil()
		}
	}

	// Parse the output to get the virtualization solution name
	virtualization := strings.TrimSpace(string(output))

	return Virtualization{
		Type: virtualization,
	}, a.ErrorOrNil()
}
