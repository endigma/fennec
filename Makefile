linux:
	@GOOS=linux GOARCH=amd64 /usr/bin/go build -tags netgo -a -v -o bin/watcher_linux_amd64
	
compressed: linux
	@upx bin/watcher_linux_amd64

run:
	@/usr/bin/go build -o bin/watcher_linux_amd64
	@bin/watcher_linux_amd64