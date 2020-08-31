OSList = darwin linux windows
ARCHList = amd64
default:
	@echo "Build executables..."
	@for os in $(OSList) ; do            \
		for arch in $(ARCHList) ; do     \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -a -ldflags '-extldflags "-static"' \
			-ldflags="-s -w" -o builds/keepmebot-$$os-$$arch .;    \
		done                              \
	done


clean:
	@rm -rf builds
	@rm -f test.db

test:
	@go test -v -cover
	@rm -f test.db

