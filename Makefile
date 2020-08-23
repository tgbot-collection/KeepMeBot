OS = darwin linux windows
ARCH = amd64
default:
	@echo "Build executables..."
	@for o in $(OS) ; do            \
		for a in $(ARCH) ; do     \
			CGO_ENABLED=0 GOOS=$$o GOARCH=$$a go build -ldflags="-s -w" -o builds/keepmebot-$$o-$$a .;    \
		done                              \
	done
	# cd builds;upx *
#	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o build/blog-linux-arm64 main.go;
#	@CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o build/blog-linux-arm main.go;


clean:
	@rm -rf builds

tests:
	@go test -v ./test/...

