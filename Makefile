OSList = darwin linux windows
ARCHList = amd64
default:
	bash autogen.sh
	@echo "Build executables..."
	@for os in $(OSList) ; do            \
		for arch in $(ARCHList) ; do     \
			GOOS=$$os GOARCH=$$arch go build -o builds/keepmebot-$$os-$$arch .;    \
		done                              \
	done


clean:
	@rm -rf builds
	@rm -f test.db

test:
	@go test -v -cover
	@rm -f test.db

