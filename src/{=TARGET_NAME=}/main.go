package main

import (
	_ "conf"
	"github.com/cjey/slog"
)

func main() {
	slog.Emerg("Hello world!")
	slog.Alert("Hello world!")
	slog.Crit("Hello world!")
	slog.Err("Hello world!")
	slog.Warning("Hello world!")
	slog.Notice("Hello world!")
	slog.Info("Hello world!")
	slog.Debug("Hello world!")
}
