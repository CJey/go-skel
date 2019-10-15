## Summary

Go skeleton，copy本项目后，请在项目根目录执行`./init <app name>`来初始化

初始化完成后，再将`init`文件和`.git`目录删除，即可独立

本项目致力于解决在开发一般性的服务/命令行程序过程中，频繁要处理解决的几个基础性问题

* 日志方面
    * 使用了[uber的zap](https://godoc.org/go.uber.org/zap)
    * 默认吐出json格式的日志，方便今后的日志解析
    * 没有pretty输出，有点不利于人肉翻读大量日志
* 命令行方面
    * 使用了[cobra](https://godoc.org/github.com/spf13/cobra)
    * 主要提供子命令风格的命令行参数解析，方便集成多条子命令
* 配置文件方面
    * 使用了[viper](https://godoc.org/github.com/spf13/viper)
    * 日志文件默认使用TOML格式
    * 将日志文件内部的参数进行半自动绑定
* 版本号方面
    * go在编译时，自动提取git信息，并注入到内建的build包中
    * 用于解决数字版本号不能准确反应代码变更的问题
    * 以及自动化的附着和变更版本号
* 编译和开发方面
    * 由于go的编译速度较快，因此可以方便的做到代码变更后快速重新编译重新运行
    * 为了让编译出的程序可以更快速的被执行，被调试，只需将main package的名称`<appname>`在项目根目录下软连接到`utils/build`上
    * 日常开发时，代码变更保存后，直接执行项目根目录下的app软链`./<appname>`就能自动编译运行

## 开发

### 环境依赖

* gcc
* golang >= 1.11
* 科学网络，顺畅下载依赖包

### 开发调试

``` shell
./{=APPNAME=}
```

例.

``` shell
./{=APPNAME=} version

# 输出
Appname      {=APPNAME=}
Version      0.0.0

GitTrace     5.8df6caf
GitBranch    master
GitHash      8df6caffcf4197aed825c3cc39ec4e66e79162da @ 2018-11-15 22:03:52
GitRepo      git@github.com:CJey/go-skel.git

BuildHash    9b50ee3834e43d79ebd20d17644c4a24ffe161a4
BuildInfo    go-1.11.2-linux/amd64 @ 2018-11-15 22:04:02
```

### 编译

``` shell
./utils/build {=APPNAME=}
```
