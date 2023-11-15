# Makefile for festivals-server-tools

VERSION=development
DATE=$(shell date +"%d-%m-%Y-%H-%M")
REF=refs/tags/development
export

build:
	go build

test:
	 go test -v

coverage:
	 go test -coverprofile=coverage_report.out
	 go tool cover -html=coverage_report.out