all: test win64 linux64

BIN=${GOPATH}/bin

test:
	go test github.com/breed85/gojenkins/slave

win64:
	mkdir -p $(BIN)
	GOOS=windows GOARCH=amd64 go build -o $(BIN)/gojenkins-windows-x64.exe

linux64:
	mkdir -p $(BIN)
	GOOS=linux GOARCH=amd64 go build -o $(BIN)/gojenkins-linux-x64
