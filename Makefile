all: win linux

WIN64=${GOPATH}/bin/windows/amd64
LINUX64=${GOPATH}/bin/linux/amd64

win:
	mkdir -p $(WIN64)
	GOOS=windows GOARCH=amd64 go build -o $(WIN64)/gojenkins.exe

linux:
	mkdir -p $(LINUX64)
	GOOS=linux GOARCH=amd64 go build -o $(LINUX64)/gojenkins
