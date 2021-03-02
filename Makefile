linux:
	@GOOS=linux GOARCH=amd64 /usr/bin/go build -tags netgo -a -v -o bin/fennec_linux_amd64
	
upx: linux
	@upx bin/fennec_linux_amd64

run:
	@/usr/bin/go build -o bin/fennec_linux_amd64
	@bin/fennec_linux_amd64