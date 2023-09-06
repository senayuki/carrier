build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -ldflags="-s -w" -o carrier-linux-amd64
	upx -9 carrier-linux-amd64