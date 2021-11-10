package version

import "fmt"

const (
	githash = "2021-11-09 18:02:16 +0800 @8880f4b1f75e08ad4c31b850f20e1e7bcc980fd4"
    branch = "master"
	buildat = "2021-11-10 16:36:13 +0800 by go version go1.17.2 windows/amd64"
    host = "rui"
)

func Show() string {
	return fmt.Sprintf("git: %v\nbranch: %v\nbuild: %v\nhost: %v", githash, branch, buildat, host)
}


