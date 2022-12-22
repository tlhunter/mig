build:
	go build

release:
	go build -ldflags="-s -w"

tiny:
	go build -ldflags="-s -w"
	upx mig
