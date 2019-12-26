package app

import (
	"fmt"
)

func Example() {
	appname = "go-skel"
	version = "1.0.0"

	gitTag = "v1"
	gitRepo = "https://github.com/cjey/go-skel"
	gitBranch = "master"
	gitHash = "dc107ed483a0b3926a357c552bd3055edb2e3207"
	gitTime = "1572676521"
	gitNumber = "1024"
	gitStatusNumber = "3"
	gitStatusHash = "212e29f41217dcb9ced3400d8f43e01c41e68d17"

	buildID = "aecb18f37b0c0051982caa2a5a42946b59e95cf7"
	buildTime = "1572676521.153552300"
	goVersion = "go version go1.13.1 linux/amd64"

	collectInfo(&App)

	fmt.Printf("%s", Info())
	// Output:
	// AppName     go-skel
	// Version     1.0.0
	//
	// GitTrace    1024.dc107ed+3.212e29f
	// GitTag      v1
	// GitBranch   master
	// GitRepo     https://github.com/cjey/go-skel
	// GitHash     dc107ed483a0b3926a357c552bd3055edb2e3207 @ 2019-11-02 14:35:21
	//
	// Golang      1.13.1 linux/amd64
	// BuildInfo   aecb18f37b0c0051982caa2a5a42946b59e95cf7 @ 2019-11-02 14:35:21
}
