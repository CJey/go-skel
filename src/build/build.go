package build

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	NOT_AVAILABLE = "N/A"
)

var (
	appname   string
	version   string
	goVersion string
	codeRoot  string

	gitRepo         string
	gitBranch       string
	gitHash         string
	gitNumber       string
	gitStatusNumber string
	gitStatusHash   string

	buildRand      string
	buildIndicator string
	buildTime      string
)

func Info() string {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "Appname      %s\n", Appname())
	fmt.Fprintf(buf, "Version      %s\n", Version())
	fmt.Fprintf(buf, "\n")
	if GitTrace() != NOT_AVAILABLE {
		fmt.Fprintf(buf, "GitTrace     %s\n", GitTrace())
		fmt.Fprintf(buf, "GitBranch    %s\n", GitBranch())
		fmt.Fprintf(buf, "GitHash      %s @ %s\n", GitHash(), GitTimeString())
		fmt.Fprintf(buf, "GitRepo      %s\n", GitRepo())
		fmt.Fprintf(buf, "\n")
	}
	fmt.Fprintf(buf, "BuildHash    %s\n", BuildHash())
	fmt.Fprintf(buf, "BuildInfo    go-%s-%s @ %s\n",
		GoVersion(), GoArch(), BuildTimeString(),
	)
	return buf.String()
}

func Appname() string {
	if len(appname) > 0 {
		return appname
	}
	return NOT_AVAILABLE
}

func Version() string {
	if len(version) > 0 {
		if version[0] != 'v' {
			return version
		}
		if len(version) > 1 {
			return version[1:]
		}
	}
	if len(appname) > 0 {
		return "0.0.0"
	}
	return NOT_AVAILABLE
}

func GitRepo() string {
	if len(gitHash) >= 40 && len(gitRepo) > 0 {
		return gitRepo
	}
	return NOT_AVAILABLE
}

func GitBranch() string {
	if len(gitHash) >= 40 && len(gitBranch) > 0 {
		return gitBranch
	}
	return NOT_AVAILABLE
}

func GitTrace() string {
	if GitNumber() > 0 {
		var devFlag string
		if len(gitStatusNumber) > 0 && len(gitStatusHash) >= 40 &&
			gitStatusHash != "da39a3ee5e6b4b0d3255bfef95601890afd80709" { // empty string sha1sum
			devFlag = fmt.Sprintf(" # %s.%s", gitStatusNumber, gitStatusHash[:7])
		}
		return fmt.Sprintf("%d.%s%s", GitNumber(), GitShortHash(), devFlag)
	}
	return NOT_AVAILABLE
}

func GitHash() string {
	if len(gitHash) >= 40 {
		return gitHash[:40]
	}
	return NOT_AVAILABLE
}

func GitTime() time.Time {
	if len(gitHash) > 41 {
		t, err := strconv.ParseInt(gitHash[41:], 10, 64)
		if err == nil {
			return time.Unix(t, 0)
		}
	}
	return time.Time{}
}

func GitTimeString() string {
	t := GitTime()
	if !t.IsZero() {
		return t.Format("2006-01-02 15:04:05")
	}
	return NOT_AVAILABLE
}

func GitShortHash() string {
	if len(gitHash) >= 40 {
		return gitHash[:7]
	}
	return NOT_AVAILABLE
}

func GitNumber() uint64 {
	n, err := strconv.ParseUint(gitNumber, 10, 64)
	if err == nil {
		return n
	}
	return 0
}

func GoVersion() string {
	tmp := strings.Split(goVersion, " ")
	if len(tmp) == 4 && len(tmp[2]) > 2 {
		return tmp[2][2:]
	}
	return NOT_AVAILABLE
}

func GoArch() string {
	tmp := strings.Split(goVersion, " ")
	if len(tmp) == 4 {
		return tmp[3]
	}
	return NOT_AVAILABLE
}

func BuildHash() string {
	if len(buildRand) > 0 {
		raw := strings.Join([]string{
			appname,
			version,
			goVersion,
			codeRoot,

			gitRepo,
			gitBranch,
			gitHash,
			gitNumber,
			gitStatusNumber,
			gitStatusHash,

			buildRand,
			buildIndicator,
			buildTime,
		}, "\x00")
		return fmt.Sprintf("%x", sha1.Sum([]byte(raw)))
	}
	return NOT_AVAILABLE
}

func BuildTime() time.Time {
	t, err := strconv.ParseInt(buildTime, 10, 64)
	if err == nil {
		return time.Unix(t, 0)
	}
	return time.Time{}
}

func BuildTimeString() string {
	t := BuildTime()
	if !t.IsZero() {
		return t.Format("2006-01-02 15:04:05")
	}
	return NOT_AVAILABLE
}

func BuildIndicator() string {
	if len(buildIndicator) > 0 {
		return buildIndicator
	}
	return NOT_AVAILABLE
}

func CodeRoot() string {
	return codeRoot
}
