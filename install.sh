#!/usr/bin/env bash

: << !
Name: install.sh
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-09-19 11:41:03

Description: 一键安装 curator
参考了：https://github.com/liuchengxu/space-vim

Attentions:
-

Depends:
- curl
- jq
!

####################################################################
#+++++++++++++++++++++++++ Define Variable ++++++++++++++++++++++++#
####################################################################
#------------------------- Program Variable
# 仓库所有人
OWNER="YHYJ"
# 仓库名
REPO="curator"

# 系统信息
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
ATYPE="tar.gz"

# 最新 Release 信息 API
RELEASE_API="https://api.github.com/repos/$OWNER/$REPO/releases/latest"
# 请求头参数
ACCEPT="application/vnd.github+json"

# 下载文件存储地址
TEMP_DIR="/tmp/manager/release/$REPO"
# 下载文件解压地址
CACHE_DIR="/tmp/manager/release/$REPO/cache"

# 程序文件信息
PROGRAM_FILE="curator"
PROGRAM_MODE="0755"
PROGRAM_OWNER="root"
PROGRAM_GROUP="root"
PROGRAM_FILE_PATH="/usr/local/bin"
# LICENSE 文件信息
LICENSE_FILE="LICENSE"
LICENSE_MODE="0644"
LICENSE_OWNER="root"
LICENSE_GROUP="root"
LICENSE_FILE_PATH="/usr/local/share/licenses/$PROGRAM_FILE"

####################################################################
#+++++++++++++++++++++++++ Define Function ++++++++++++++++++++++++#
####################################################################
#------------------------- Exit Code Variable
readonly err_file=1         # 文件/路径类错误
readonly err_no_program=127 # 未找到命令
readonly err_unsupport=202  # 不受支持
#------------------------- Feature Function
function msg() { # 输出基础信息
  echo -e "$1" >&2
}

function success() { # 输出成功信息
  msg "\x1b[32m[✔]\x1b[0m ${1}"
}

function error() { # 输出错误信息
  msg "\x1b[31m[✘]\x1b[0m ${1}"
  exit "${2}"
}

function checkos() { # 检查系统
  # 系统类型
  case $OS in
    linux)
      ATYPE="tar.gz"
      ;;
    darwin)
      ATYPE="tar.gz"
      ;;
    *)
      # 不支持其他系统
      error "Unsupported system: $OS" $err_unsupport
      ;;
  esac

  # 系统架构
  case $ARCH in
    x86_64)
      ARCH="amd64"
      ;;
    aarch64)
      ARCH="arm64"
      ;;
    armv7l)
      ARCH="armv7"
      ;;
    armv6l)
      ARCH="armv6"
      ;;
    *)
      # 如果不是已知的架构，保留原始值
      ;;
  esac
}

function exists() { # 检查命令是否存在
  command -v "$1" > /dev/null 2>&1
}

function checkdep() { # 检查依赖
  if ! exists "$1"; then
    error "Missing dependency: $1" $err_no_program
  fi
}

function checkpath() { # 检查文件夹是否存在
  if [ ! -d "$1" ]; then
    mkdir -p "$1"
  fi
}

function goto() { # 跳转到指定目录
  checkpath "$1"
  cd "$1" || exit
}

function download() { # 下载文件（$1 是下载链接，$2 是下载的文件名, $3 是下载文件大小）
  if [ -n "$1" ]; then
    if exists "pv" && [[ $3 -gt 0 ]]; then
      # 使用 pv 的情况下写入文件由 pv 控制，curl 不使用 '-O' 参数
      curl -L -s "$1" | pv -s "$3" -N "Downloading $2" > "$2"
    else
      curl -LO -s "$1"
    fi
  else
    error "No download link provided" $err_file
  fi
}

function checksum() { # 校验文件（不校验找不到的文件；仅使用状态码表示结果）
  if [ -n "$1" ]; then
    if ! sha256sum --ignore-missing --status --check "$1"; then
      error "File verification failed" $err_file
    fi
  fi
}

function installer() { # 安装程序
  checkpath "$PROGRAM_FILE_PATH"
  install --mode="$PROGRAM_MODE" --owner="$PROGRAM_OWNER" --group="$PROGRAM_GROUP" "$CACHE_DIR/$PROGRAM_FILE" "$PROGRAM_FILE_PATH/$PROGRAM_FILE"
  checkpath "$LICENSE_FILE_PATH"
  install --mode="$LICENSE_MODE" --owner="$LICENSE_OWNER" --group="$LICENSE_GROUP" "$CACHE_DIR/$LICENSE_FILE" "$LICENSE_FILE_PATH/$LICENSE_FILE"
}

####################################################################
#++++++++++++++++++++++++++++++ Main ++++++++++++++++++++++++++++++#
####################################################################
# 检查系统
checkos
# 检查依赖项
checkdep 'curl'
checkdep 'jq'

# 访问 GitHub API，获取 Release 信息
response=$(curl -L -H "Accept: $ACCEPT" -s "$RELEASE_API")
assets=$(echo "$response" | jq -r '.assets[] | {name: .name, size: .size, content_type: .content_type, download_url: .browser_download_url, download_count: .download_count}')

# 获取 TAG
TAG=$(echo "$response" | jq -r '.tag_name')
# 校验文件名
CHECKSUMS_FILE="checksums.txt"
# 程序压缩文件名
ARCHIVE_FILE="${REPO}_${TAG}_${OS}_${ARCH}.${ATYPE}"

# 获取下载文件信息
checksums_file_download_url=$(echo "$assets" | jq -r 'select(.name=="'"$CHECKSUMS_FILE"'") | .download_url')
checksums_file_size=$(echo "$assets" | jq -r 'select(.name=="'"$CHECKSUMS_FILE"'") | .size')
archive_file_download_url=$(echo "$assets" | jq -r 'select(.name=="'"$ARCHIVE_FILE"'") | .download_url')
archive_file_size=$(echo "$assets" | jq -r 'select(.name=="'"$ARCHIVE_FILE"'") | .size')

# 跳转到临时目录
goto "$TEMP_DIR"

# 下载文件
download "$checksums_file_download_url" "$CHECKSUMS_FILE" "$checksums_file_size"
download "$archive_file_download_url" "$ARCHIVE_FILE" "$archive_file_size"

# 校验文件
checksum "$CHECKSUMS_FILE"

# 解压
checkpath "$CACHE_DIR"
tar -xzpf "$ARCHIVE_FILE" --directory="$CACHE_DIR"

# 安装
installer
success "Successfully installed \x1b[32m$REPO\x1b[0m"

# 清理垃圾
rm -rf "$TEMP_DIR"
