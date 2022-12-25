build:
	go build

release:
	go build -ldflags="-s -w" -o mig

tiny:
	go build -ldflags="-s -w"
	upx mig

multi:
	echo "linux amd64"
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o mig-linux-amd64
	echo "windows amd64"
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o mig.exe
	echo "macos amd64"
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o mig-macos-amd64
	echo "macos arm64"
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o mig-macos-arm64
