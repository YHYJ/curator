# README

<!-- File: README.md -->
<!-- Author: YJ -->
<!-- Email: yj1516268@outlook.com -->
<!-- Created Time: 2023-04-18 13:19:11 -->

---

## Table of Contents

<!-- vim-markdown-toc GFM -->

* [示例配置](#示例配置)

<!-- vim-markdown-toc -->

---

<!-- Object info -->

---

用于克隆指定用户的指定仓库

## 示例配置

```toml
[ssh]
private_key_file = "/home/yj/.ssh/id_rsa"

[storage]
path = "/home/yj/Documents/Repos"

[git]
url = "git@github.com:YHYJ"
repos = ["Repo_1", "Repo_2"]
```
