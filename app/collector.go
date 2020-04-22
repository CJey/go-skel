package app

import (
	"encoding/base64"
	"strconv"
	"strings"
	"time"
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

	// fullversion
	if len(git.CommitHash) >= 7 {
		app.FullVersion = app.Version + "-" + strconv.Itoa(int(app.Release)) + "." + git.CommitHash[:7]
	} else {
		app.FullVersion = app.Version + "-" + strconv.Itoa(int(app.Release))
	}
}

func base64d(enc string) string {
	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func beUint(s string) uint {
	if len(s) == 0 {
		return 0
	}
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return uint(n)
}

func beTrace(n uint, h string) string {
	if n > 0 && len(h) >= 7 {
		return strconv.Itoa(int(n)) + "." + h[:7]
	}
	return ""
}

func beTime(s string) time.Time {
	if len(s) == 0 {
		return time.Time{}
	}

	var sec, nsec int64
	var err error
	// sec[.nsec] <=> 1572017404[.638981238]
	ss := strings.SplitN(s, ".", 2)
	if len(ss[0]) > 0 {
		sec, err = strconv.ParseInt(ss[0], 10, 64)
		if err != nil {
			panic(err)
		}
	}
	if len(ss) > 1 && len(ss[1]) > 0 {
		nsec, _ = strconv.ParseInt(ss[1], 10, 64)
	}

	return time.Unix(sec, nsec)
}
