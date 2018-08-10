.PHONY: gher gher-linux-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

gher: 
	go run build/ci.go install
	@echo "Compilation done"
	echo "Run \"$(GOBIN)/gher\" to launch gher."
	
gher-linux-amd64:
	go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v
	@echo "Linux amd64 cross compilation done"
	
clean:
	rm -fr $(GOBIN)/*