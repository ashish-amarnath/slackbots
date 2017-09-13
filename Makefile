image_name := slackbot


.PHONY: app test clean pkg pkg_test

mkbuilddir:
	mkdir -p build

pkg:
	go build pkg/types/*.go
	go build pkg/slack/*.go

pkg-test: pkg
	go test pkg/...

cmd:
	go build cmd/*.go

cmd-test: cmd
	go test cmd/...

app:
	go build -o build/slackbot *.go

