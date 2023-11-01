PROJECT       := github.com/yhyj/clone-repos
TARGET        := clone-repos
GENERATE_PATH := build
INSTALL_PATH  := /usr/local/bin
ATTRIBUTION   := root
COMMIT        := $(shell git rev-parse HEAD)

.PHONY: all tidy build install clean
all: build

help:
	@echo "usage: make [OPTIONS]"
	@echo "    help        Show this message"
	@echo "    tidy        Update project module dependencies"
	@echo "    build       Compile and generate executable file"
	@echo "    install     Install executable file"
	@echo "    clean       Clean build process files"

tidy:
	@go mod tidy
	@echo -e "\x1b[32m[✔]\x1b[0m Successfully tidied up dependencies"

build:
	@go build -gcflags="-trimpath" -ldflags="-s -w -X $(PROJECT)/general.GitCommitHash=$(COMMIT) -X $(PROJECT)/general.BuildTime=`date +%s` -X $(PROJECT)/general.BuildBy=Makefile" -o $(GENERATE_PATH)/$(TARGET)
	@echo -en "\x1b[32m[✔]\x1b[0m Successfully generated \x1b[32;1m$(TARGET)\x1b[0m"

install:
	@install --mode=755 --owner=$(ATTRIBUTION) --group=$(ATTRIBUTION) $(GENERATE_PATH)/$(TARGET) $(INSTALL_PATH)/$(TARGET)
	@echo -e "\r\x1b[K\x1b[0m\x1b[32m[✔]\x1b[0m Successfully installed \x1b[32m$(TARGET)\x1b[0m"

clean:
	@rm -rf $(GENERATE_PATH) && echo -e "    - Removed \x1b[32m$(GENERATE_PATH)\x1b[0m"
	@echo -e "\x1b[32m[✔]\x1b[0m Successfully cleared build folder"
