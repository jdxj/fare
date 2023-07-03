OUTPUT := output

.PHONY: build.%
build.%:
	@mkdir -p $(OUTPUT)
	@CGO_ENABLED=0 GOOS=$* GOARCH=amd64 go build -ldflags '-s -w' -o $(OUTPUT)/fare_$* cmd/*.go

.PHONY: clear
clear:
	@rm -rvf $(OUTPUT)
