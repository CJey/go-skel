## Summary

Go skeleton

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
