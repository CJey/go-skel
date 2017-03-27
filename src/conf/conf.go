package conf

import (
	"os"

	"version"

	"github.com/alecthomas/kingpin"
)

var showVersion bool

var (
	UseSyslog   bool
	LogLevel    string
	LogLineOff  bool
	LogLevelOff bool
	LogTimeOff  bool
)

func parseFlags() {
	kingpin.Flag("version", "Show Version Info").Default("false").BoolVar(&showVersion)
	kingpin.Flag("syslog", "Output redirect to syslog").Default("false").BoolVar(&UseSyslog)
	kingpin.Flag("log-level", "log level").
		Default("info").
		EnumVar(&LogLevel, "emerg", "alert", "crit", "err", "warning", "notice", "info", "debug")
	kingpin.Flag("log-lineoff", "Hide the code line of log").Default("false").BoolVar(&LogLineOff)
	kingpin.Flag("log-leveloff", "Hide the log level hint string").Default("false").BoolVar(&LogLevelOff)
	kingpin.Flag("log-timeoff", "Hide the time of log").Default("false").BoolVar(&LogTimeOff)

	switch version.Name() {
	case "Server":
		parseServer()
	case "Client":
		parseClient()
	}

	kingpin.Parse()
}

func prepareFlags() {
	initLog()

	switch version.Name() {
	case "Server":
		prepareServer()
	case "Client":
		prepareClient()
	}
}

// Server flags

var (
// server args
)

func parseServer() {
}

func prepareServer() {
}

// Client flags

var (
// client args
)

func parseClient() {
}

func prepareClient() {
}

func init() {
	parseFlags()

	if showVersion {
		print(version.Show())
		os.Exit(0)
	}

	prepareFlags()
}
