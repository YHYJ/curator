#!/usr/bin/env bash

: << !
Name: create-hook-link.sh
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-08 14:14:21

Description: 创建git钩子

Attentions:
-

Depends:
-
!

mkdir .git/hooks
ln -sf ../../hooks/post-commit .git/hooks/post-commit
