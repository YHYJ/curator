# README

<!-- File: README.md -->
<!-- Author: YJ -->
<!-- Email: yj1516268@outlook.com -->
<!-- Created Time: 2023-04-18 13:19:11 -->

---

## Table of Contents

<!-- vim-markdown-toc GFM -->

* [Usage](#usage)
* [Compile](#compile)
  * [当前平台](#当前平台)
  * [交叉编译](#交叉编译)
    * [Linux](#linux)
    * [macOS](#macos)
    * [Windows](#windows)

<!-- vim-markdown-toc -->

---

<!-------------------------------------------------------------->
<!--       _                                                  -->
<!--   ___| | ___  _ __   ___       _ __ ___ _ __   ___  ___  -->
<!--  / __| |/ _ \| '_ \ / _ \_____| '__/ _ \ '_ \ / _ \/ __| -->
<!-- | (__| | (_) | | | |  __/_____| | |  __/ |_) | (_) \__ \ -->
<!--  \___|_|\___/|_| |_|\___|     |_|  \___| .__/ \___/|___/ -->
<!--                                        |_|               -->
<!-------------------------------------------------------------->

---

用于克隆指定用户的指定仓库

## Usage

- `config`子命令

    该子命令用于操作配置文件，有以下参数：

    - 'create'：创建默认内容的配置文件，可以使用全局参数'--config'指定配置文件路径
    - 'force'：当指定的配置文件已存在时，使用该参数强制覆盖原文件
    - 'print'：打印配置文件内容

- `run`子命令

    使用该子命令开始执行克隆

- `version`子命令

    查看程序版本信息

- `help`

    查看程序帮助信息

## Compile

### 当前平台

```bash
go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/clone-repos/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/clone-repos/general.BuildTime=`date +%s` -X github.com/yhyj/clone-repos/general.BuildBy=$USER" -o build/clone-repos main.go
```

### 交叉编译

使用命令`go tool dist list`查看支持的平台

#### Linux

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/clone-repos/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/clone-repos/general.BuildTime=`date +%s` -X github.com/yhyj/clone-repos/general.BuildBy=$USER" -o build/clone-repos main.go
```

> 使用`uname -m`确定硬件架构
>
> - 结果是x86_64则GOARCH=amd64
> - 结果是aarch64则GOARCH=arm64

#### macOS

```bash
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/clone-repos/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/clone-repos/general.BuildTime=`date +%s` -X github.com/yhyj/clone-repos/general.BuildBy=$USER" -o build/clone-repos main.go
```

> 使用`uname -m`确定硬件架构
>
> - 结果是x86_64则GOARCH=amd64
> - 结果是aarch64则GOARCH=arm64

#### Windows

```powershell
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -gcflags="-trimpath" -ldflags="-s -w -H windowsgui -X github.com/yhyj/clone-repos/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/clone-repos/general.BuildTime=`date +%s` -X github.com/yhyj/clone-repos/general.BuildBy=$USER" -o build/clone-repos.exe main.go
```

> 使用`echo %PROCESSOR_ARCHITECTURE%`确定硬件架构
>
> - 结果是x86_64则GOARCH=amd64
> - 结果是aarch64则GOARCH=arm64
