BuildStamp = main.BuildStamp=$(shell date '+%Y-%m-%d_%I:%M:%S%p')
GitHash    = main.GitHash=$(shell git rev-parse HEAD)
Version    = main.Version=$(shell git describe --abbrev=0 --tags --always)
Target     = $(shell basename $(abspath $(dir $$PWD)))

all: $(shell uname -s)

Darwin:
	GOOS=darwin GOARCH=amd64 go build -v -o release/${Target}-drawin -ldflags "-s -w -X ${BuildStamp} -X ${GitHash} -X ${Version}"

Linux:
	GOOS=linux GOARCH=amd64 go build -v -o release/${Target}-linux -ldflags "-s -w -X ${BuildStamp} -X ${GitHash} -X ${Version}"

armv7l:
	GOOS=linux GOARCH=arm GOARM=7 go build -v -o release/${Target}-armv7 -ldflags "-s -w -X ${BuildStamp} -X ${GitHash} -X ${Version}"

authors:
	printf "Authors\n=======\n\nProject's contributors:\n\n" > AUTHORS.md
	git log --raw | grep "^Author: " | cut -d ' ' -f2- | cut -d '<' -f1 | sed 's/^/- /' | sort | uniq >> AUTHORS.md

clean:
	-rm -rf release *.db *.db-journal
