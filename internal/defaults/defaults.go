package defaults

import "os"

const (
	DEFAULT_ENV_PATH = "/etc/env/.certwatch.env"
)

func MachineName() string {
	machineName, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return machineName
}
