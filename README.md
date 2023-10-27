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

- 编译当前平台可执行文件：

```bash
go build main.go
```

- **交叉编译**指定平台可执行文件：

```bash
# 适用于Linux AArch64平台
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build main.go
```

```bash
# 适用于macOS amd64平台
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
```

```bash
# 适用于Windows amd64平台
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
```
