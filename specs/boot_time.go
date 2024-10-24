package specs

import (
	"time"

	"github.com/shirou/gopsutil/host"
)

func fetchBootTime() (time.Time, error) {
	bootTime, err := host.BootTime()
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(int64(bootTime), 0), nil
}
