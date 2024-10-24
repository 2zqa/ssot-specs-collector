package specs

import (
	"os/user"
)

func hasNoRoot() bool {
	current, err := user.Current()
	if err != nil {
		return true
	}

	return current.Uid != "0"
}
