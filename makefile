# -------- CONFIG --------
APP_NAME := taskcli
VERSION  := 1.0.0
ARCH     := amd64

PKG_DIR     := packaging/deb
BIN_DIR     := $(PKG_DIR)/usr/local/bin
SERVICE_DIR := $(PKG_DIR)/etc/systemd/system
CONTROL_DIR := $(PKG_DIR)/DEBIAN

OUTPUT_DEB := $(APP_NAME)_$(VERSION)_$(ARCH).deb

# -------- PHONY TARGETS --------
.PHONY: all build package clean install

# -------- DEFAULT TARGET --------
all: package

# -------- BUILD BINARY --------
build:
	go build -o $(APP_NAME)

# -------- PACKAGE AS .DEB --------
package: build
	mkdir -p $(BIN_DIR)
	mkdir -p $(SERVICE_DIR)
	mkdir -p $(CONTROL_DIR)
	cp $(APP_NAME) $(BIN_DIR)/
	cp packaging/linux/$(APP_NAME).service $(SERVICE_DIR)/
	cp packaging/linux/control $(CONTROL_DIR)/control
	dpkg-deb --build $(PKG_DIR) $(OUTPUT_DEB)

# -------- CLEAN EVERYTHING --------
clean:
	rm -f $(APP_NAME)
	rm -f $(OUTPUT_DEB)
	rm -rf $(PKG_DIR)/usr
	rm -rf $(PKG_DIR)/etc
	rm -rf $(PKG_DIR)/DEBIAN

# -------- INSTALL PACKAGE --------
install: package
	sudo dpkg -i $(OUTPUT_DEB)
