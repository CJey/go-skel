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
	_App Application

	ID          string // _App.Build.ID
	Name        string // _App.Name
	Trace       string // _App.Git.Trace
	Version     string // _App.Version
	Release     uint   // _App.Release
	FullVersion string // _App.FullVersion
)

func init() {
	collectInfo(&_App)

	ID = _App.Build.ID
	Name = _App.Name
	Trace = _App.Git.Trace
	Version = _App.Version
	Release = _App.Release
	FullVersion = _App.FullVersion
}

func App() Application {
	return _App
}

// Application，描述应用程序编译时的环境信息
type Application struct {
	Name        string // 名称
	Version     string // 语义化版本
	Release     uint   // 第几次发布
	FullVersion string // 完整版本

	Git struct {
		Repo   string // 编译分支所track的upstream地址
		Branch string // 编译的分支名称，如果为HEAD，则会被置空
		// 结合CommitTrace和StatusTrace生成的唯一标记
		// 一般情况下，如果Trace相同，则意味着两个binary在功能上是等价的
		// 只不过是从不同的渠道编译而来
		Trace string

		CommitHash       string    // 编译时所在commit的hash
		CommitNumber     uint      // 累加计算，得到本commit从第一次commit开始至今的数量值
		CommitTime       time.Time // 编译时所在commit的生成时间
		CommitTimeString string    // 时间的快捷字符串标记，方便template引用
		CommitTrace      string    // 结合number和hash生成的一组标记，标记此commit的唯一性

		TagName       string // 编译时git describe计算得到的tag名称，可能为空
		TagHash       string
		TagNumber     uint
		TagTime       time.Time
		TagTimeString string
		TagTrace      string
		TagDiff       uint // 此tag距离当前的commit相距几个commit
		TagMessage    string

		// 编译时，工作目录可能并不干净，尤其是开发者在开发过程中
		// 做的临时编译结果，因此，需要一定的手段来帮助标记此情况
		// git status的结果可以用于一定的参考
		// 如果StatusTrace不为空，则基本意味着，本程序来自于有未commit代码的目录下的编译结果
		// 一定程度上意味着这是一个临时测试版本
		StatusHash   string // git status相关的所有文件内容的hash结果
		StatusNumber uint   // git status相关的所有文件总数
		StatusTrace  string // 结合number和hash生成的一组标记，标记此工作目录的唯一性
	}

	Go struct {
		Arch    string // 编译所用golang的架构
		Version string // 编译所用golang的版本号
	}

	Build struct {
		// 每一次编译，都由编译工具生成一个随机串，标记本次编译的唯一性
		// 一般情况下，ID相同的binary文件就是同一个文件
		ID         string
		Time       time.Time // 编译时的时间
		TimeString string    // 时间的快捷字符串标记，方便template引用
		Magic      string    // 编译工具任意注入的字符串，主要用于在开发过程中由编译工具参考是否需要重新编译
	}
}
