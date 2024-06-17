#!/usr/bin/env bash

: << !
Name: sync.sh
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-16 16:51:07

Description: 从 System 同步脚本和 git hook 等

Attentions:
-

Depends:
-
!

script_dir=$(dirname "$0") # 本脚本所在路径

repo_dir="${HOME}/Documents/Repos" # 本地源库路径
tmp_dir="/tmp/sync-hook"           # 临时源库路径

source_dir="System/git"      # 来源路径
save_dir="${script_dir}/git" # 本库路径

repo_name="git@github.com:YHYJ/System.git" # 云端 System 存储库

mkdir -p "${save_dir}"

if [ -d "${repo_dir}/${source_dir}" ]; then
  echo "Synchronizing hooks from local 'System' repo"
  cp -r "${repo_dir}/${source_dir}/"* "${save_dir}"
  echo "Hooks copied successfully"
else
  echo "Synchronizing hooks from cloud 'System' repo"
  git clone "${repo_name}" "${tmp_dir}/System"
  cp -r "${tmp_dir}/${source_dir}/"* "${save_dir}"
  rm -rf "${tmp_dir}"
  echo "Hooks copied successfully"
fi

cp -r "${save_dir}/scripts/"* "${PWD}"
rm -rf "${save_dir}/scripts"
