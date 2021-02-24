linux:
	@GOOS=linux GOARCH=amd64 /usr/bin/go build -tags netgo -a -v -o bin/watcher_linux_amd64
	
compress:
	@upx bin/watcher_linux_amd64
