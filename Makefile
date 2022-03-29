.PHONY: update clean build build-all run package deploy test authors dist

NAME 					:= geoip
VERSION                 := $(shell cat VERSION)
LDFLAGS                 := -ldflags "-w -s --extldflags '-static'"

default: linux

linux: build-eduid_ladok-linux
mac_m1: build-eduid_ladok-arm 

build-eduid_ladok-arm:
		@echo building ${NAME} for darwin on arm64
		@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -v -o ./bin/${NAME} ${LDFLAGS} ./cmd/main.go 
		@echo Done

build-eduid_ladok-linux:
		@echo building ${NAME} for linux on amd64
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/${NAME} ${LDFLAGS} ./cmd/main.go
		@echo Done
