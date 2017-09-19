image_name := slackbot


.PHONY: mkbuilddir buildbot runtests clean

mkbuilddir:
	mkdir -p build

buildbot: mkbuilddir
	go build -o build/kube2iam-bot ./

runtests:
	go test -v ./...

clean:
	rm -rf ./build
