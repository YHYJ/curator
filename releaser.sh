#!/usr/bin/env bash

: << !
Name: x.sh
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-06-14 16:33:16

Description: 一键释出

Attentions:
-

Depends:
-
!

# 版本文件
target_file="general/version.go"

# 版本变量
version_var="Version string"

# 获取暂存区（已暂存但未提交）文件中的版本变量值（旧）
staged_version=$(git show :$target_file | grep "$version_var" | awk -F '"' '{print $2}')
# 获取工作区（未提交的变更）文件中的版本变量值（新）
current_version=$(grep "$version_var" $target_file | awk -F '"' '{print $2}')

# 1. 判断是否是版本变量的值发生变更
if [ "$staged_version" != "$current_version" ]; then
  echo -e "\x1b[35m-----\x1b[0m \x1b[36m正在将版本更新记录到存储库\x1b[0m \x1b[35m-----\x1b[0m"
  echo -e "\x1b[33m$staged_version\x1b[0m \x1b[36m->\x1b[0m \x1b[34m$current_version\x1b[0m"
  # 2. 执行 `git add`
  git add $target_file

  echo -e ""

  # 3. 执行 `git commit`
  git commit -m "$current_version"
  # 4. 获取 commit hash
  commit_hash=$(git rev-parse --short HEAD)

  echo -e ""

  # 提示开始释出新版本
  echo -e "\x1b[35m-----\x1b[0m \x1b[36m开始释出\x1b[0m \x1b[32m$current_version\x1b[0m \x1b[35m-----\x1b[0m"
  # 5. 执行 `git tag`
  git tag "$current_version" "$commit_hash"
  # 6. 执行 `git push`
  git push
  # 7. 执行 `git push --tags`
  git push --tags
else
  echo -e "仅当 \x1b[34m$target_file\x1b[0m 文件中的 \x1b[36mVersion\x1b[0m 值有变更时可执行"
fi
