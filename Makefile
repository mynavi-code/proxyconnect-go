proxyconnect-go: main.go
	go build

proxyconnect-go.exe: main.go
	GOOS=windows GOARCH=amd64 go build

