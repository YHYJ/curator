<h1 align="center">Curator</h1>

<!-- File: README.md -->
<!-- Author: YJ -->
<!-- Email: yj1516268@outlook.com -->
<!-- Created Time: 2023-04-18 13:19:11 -->

---

<p align="center">
  <a href="https://github.com/YHYJ/curator/actions/workflows/release.yml"><img src="https://github.com/YHYJ/curator/actions/workflows/release.yml/badge.svg" alt="Go build and release by GoReleaser"></a>
</p>

---

## Table of Contents

<!-- vim-markdown-toc GFM -->

* [Install](#install)
  * [一键安装](#一键安装)
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

## Install

### 一键安装

```bash
curl -fsSL https://raw.githubusercontent.com/YHYJ/curator/main/install.sh | sudo bash -s
```

## Usage

- `config`子命令

  该子命令用于操作配置文件，有以下参数：

  - 'create'：创建默认内容的配置文件，可以使用全局参数'--config'指定配置文件路径
  - 'force'：当指定的配置文件已存在时，使用该参数强制覆盖原文件
  - 'print'：打印配置文件内容

- `clone`子命令

  使用该子命令进行克隆，有以下参数：

  - '--source'：指定使用的仓库源，目前支持 github.com 和 git.yj1516.top

- `pull`子命令

  使用该子命令拉取远端仓库最新修改，有以下参数：

  - '--source'：指定使用的仓库源，目前支持 github.com 和 git.yj1516.top

- `version`子命令

  查看程序版本信息

- `help`子命令

  查看程序帮助信息

## Compile

### 当前平台

```bash
go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/curator/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/curator/general.BuildTime=`date +%s` -X github.com/yhyj/curator/general.BuildBy=$USER" -o build/curator main.go
```

### 交叉编译

使用命令`go tool dist list`查看支持的平台

#### Linux

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/curator/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/curator/general.BuildTime=`date +%s` -X github.com/yhyj/curator/general.BuildBy=$USER" -o build/curator main.go
```

> 使用`uname -m`确定硬件架构
>
> - 结果是 x86_64 则 GOARCH=amd64
> - 结果是 aarch64 则 GOARCH=arm64

#### macOS

```bash
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/curator/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/curator/general.BuildTime=`date +%s` -X github.com/yhyj/curator/general.BuildBy=$USER" -o build/curator main.go
```

> 使用`uname -m`确定硬件架构
>
> - 结果是 x86_64 则 GOARCH=amd64
> - 结果是 aarch64 则 GOARCH=arm64

#### Windows

```powershell
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -gcflags="-trimpath" -ldflags="-s -w -H windowsgui -X github.com/yhyj/curator/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/curator/general.BuildTime=`date +%s` -X github.com/yhyj/curator/general.BuildBy=$USER" -o build/curator.exe main.go
```

> 使用`echo %PROCESSOR_ARCHITECTURE%`确定硬件架构
>
> - 结果是 x86_64 则 GOARCH=amd64
> - 结果是 aarch64 则 GOARCH=arm64
