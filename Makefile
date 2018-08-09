.PHONY: gher

GOBIN = $(shell pwd)/build/bin
GO ?= latest

gher: 
	go run build/ci.go