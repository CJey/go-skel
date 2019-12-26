// Author: cjey.hou@gmail.com

/*
配合外部工具，在编译时搜集项目的版本信息，编译环境，git信息。

一般情况下，语义化版本号的维护工作如果暂时没有精力和意愿来执行，
则可以用自动计算出的GitTrace代替版本管理&追踪功能

注意，GitTrace并不具备理论上的完美唯一性，但在实践当中基本上可以被认为是唯一的
*/
package app

import (
	"time"
)

var (
	// 默认的App，所有搜集到的编译环境信息都在其中
	app Application

	ID      string // app.Build.ID
	Name    string // app.Name
	Trace   string // app.Git.Trace
	Version string // app.Version-app.Release

	BootID   string    // app.Boot.ID
	BootTime time.Time // app.Boot.Time
)

func init() {
	collectInfo(&app)

	ID = app.Build.ID
	Name = app.Name
	Trace = app.Git.Trace
	Version = app.Version
	BootID = app.Boot.ID
	BootTime = app.Boot.Time
}

func App() Application {
	return app
}
