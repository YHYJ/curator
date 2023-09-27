TARGET       := clone-repos
INSTALL_PATH := /usr/local/bin
ATTRIBUTION  := root

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
	@echo -e "\x1b[32m==>\x1b[0m Tidying up dependencies"
	@go mod tidy

build:
	@echo -e "\x1b[32m==>\x1b[0m Trying to compile project"
	@go build -trimpath -ldflags "-s -w" -o $(TARGET)
	@echo -e "\x1b[32m[✔]\x1b[0m Successfully generated \x1b[32m$(TARGET)\x1b[0m"

install:
	@echo -e "\x1b[32m==>\x1b[0m Trying to install $(TARGET)"
	@install --mode=755 --owner=$(ATTRIBUTION) --group=$(ATTRIBUTION) $(TARGET) $(INSTALL_PATH)/$(TARGET)
	@echo -e "\x1b[32m[✔]\x1b[0m Successfully installed \x1b[32m$(TARGET)\x1b[0m"

clean:
	@echo -e "\x1b[32m==>\x1b[0m Cleaning build process files"
	@rm -f $(TARGET)     && echo -e "    - Removed \x1b[32m$(TARGET)\x1b[0m"
	@echo -e "\x1b[32m[✔]\x1b[0m Successfully cleaned files"
