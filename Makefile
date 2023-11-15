# 项目名
PROJECT := github.com/yhyj/clone-repos
# 安装文件属主/属组
ATTRIBUTION := root
# 编译结果路径
GENERATE_PATH := build
# 可执行文件名
TARGET := clone-repos
# 可执行文件安装路径
INSTALL_PATH := /usr/local/bin
# 资源文件安装路径
RESOURCE_INSTALL_PATH := /usr/local/share
# Commit哈希值
COMMIT := $(shell git rev-parse HEAD)

.PHONY: all tidy build install uninstall clean
all: build

help:
	@echo "usage: make [OPTIONS]"
	@echo "    help        Show this message"
	@echo "    tidy        Update project module dependencies"
	@echo "    build       Compile and generate executable file"
	@echo "    install     Install executable file"
	@echo "    uninstall   Uninstall executable file"
	@echo "    clean       Clean build process files"

tidy:
	@go mod tidy
	@echo -e "\x1b[32;1m[✔]\x1b[0m Successfully tidied up dependencies"

build:
	@go build -gcflags="-trimpath" -ldflags="-s -w -X $(PROJECT)/general.GitCommitHash=$(COMMIT) -X $(PROJECT)/general.BuildTime=`date +%s` -X $(PROJECT)/general.BuildBy=Makefile" -o $(GENERATE_PATH)/$(TARGET)
	@echo -en "\x1b[32;1m[✔]\x1b[0m Successfully generated \x1b[32m$(TARGET)\x1b[0m"

install:
	@install --mode=755 --owner=$(ATTRIBUTION) --group=$(ATTRIBUTION) $(GENERATE_PATH)/$(TARGET) $(INSTALL_PATH)/$(TARGET)
	@mkdir -p $(RESOURCE_INSTALL_PATH)/licenses/$(TARGET)
	@install --mode=644 --owner=$(ATTRIBUTION) --group=$(ATTRIBUTION) LICENSE $(RESOURCE_INSTALL_PATH)/licenses/$(TARGET)/LICENSE
	@echo -e "\r\x1b[K\x1b[0m\x1b[32;1m[✔]\x1b[0m Successfully installed \x1b[32m$(TARGET)\x1b[0m"

uninstall:
	@rm -rf $(INSTALL_PATH)/$(TARGET) $(RESOURCE_INSTALL_PATH)/licenses/$(TARGET)
	@echo -e "\x1b[K\x1b[0m\x1b[32;1m[✔]\x1b[0m \x1b[32m$(TARGET)\x1b[0m has been \x1b[31;1muninstalled\x1b[0m"

clean:
	@rm -rf $(GENERATE_PATH) && echo -e "    - Removed \x1b[32m$(GENERATE_PATH)\x1b[0m"
	@echo -e "\x1b[32;1m[✔]\x1b[0m Successfully cleared build folder"
