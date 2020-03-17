package app

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

// 此处定义的变量，均可能会被编译工具在编译时注入初始值
// go build -X "<module-name>/app.<global variable name>=<string value>" ...
// e.g. go build -X "go-skel/app.version=0.0.1" ...
var (
	appname string = "myapp"
	version string = "0.0.1"
	release string = "1"

	gitRepo         string
	gitBranch       string
	gitHash         string
	gitTime         string
	gitNumber       string
	gitStatusNumber string
	gitStatusHash   string

	tagName    string
	tagHash    string
	tagTime    string
	tagNumber  string
	tagDiff    string
	tagMessage string

	buildID    string
	buildTime  string
	buildMagic string
	goVersion  string
)

// 将编译工具注入的初始值解析处理后，赋值于默认的App之中
func collectInfo(app *Application) {
	const tf = "2006-01-02 15:04:05"
	// base
	app.Name = appname
	app.Version = version
	app.Release = beUint(release)

	// boot
	app.Boot.ID = uuid.NewV4().String()
	app.Boot.Time = time.Now()

	// git
	git := &app.Git
	git.Repo = gitRepo
	if gitBranch != "HEAD" { // HEAD意味着当前并不处于某个具名的分支上，但不排除正处于某个tag上
		git.Branch = gitBranch
	}

	git.CommitHash = gitHash
	git.CommitNumber = beUint(gitNumber)
	git.CommitTime = beTime(gitTime)
	git.CommitTimeString = git.CommitTime.Format(tf)
	git.CommitTrace = beTrace(git.CommitNumber, git.CommitHash)

	git.TagName = tagName
	git.TagHash = tagHash
	git.TagNumber = beUint(tagNumber)
	git.TagTime = beTime(tagTime)
	git.TagTimeString = git.TagTime.Format(tf)
	git.TagTrace = beTrace(git.TagNumber, git.TagHash)
	git.TagDiff = beUint(tagDiff)
	git.TagMessage = base64d(tagMessage)

	git.StatusHash = gitStatusHash
	git.StatusNumber = beUint(gitStatusNumber)
	git.StatusTrace = beTrace(git.StatusNumber, git.StatusHash)

	git.Trace = git.CommitTrace
	if len(git.Trace) > 0 && len(git.StatusTrace) > 0 {
		git.Trace += " + " + git.StatusTrace
	}

	// golang
	if len(goVersion) > 0 {
		ss := strings.SplitN(goVersion[13:], " ", 2)
		app.Go.Version = ss[0]
		app.Go.Arch = ss[1]
	}

	// build
	app.Build.ID = buildID
	app.Build.Time = beTime(buildTime)
	app.Build.TimeString = app.Build.Time.Format(tf)
	app.Build.Magic = buildMagic
}

func base64d(enc string) string {
	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		panic(err)
	}
	return string(data)
}
