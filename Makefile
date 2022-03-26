.PHONY: build
VERSION=`git describe --tags --abbrev=0`

build:
	goreleaser release --snapshot --skip-publish --rm-dist

release:
	goreleaser release --rm-dist
