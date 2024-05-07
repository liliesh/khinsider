package env

import (
	"os/user"
	"path/filepath"
)

func GetCachePath() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return filepath.Join(usr.HomeDir, ".cache/khinsider/")
}
