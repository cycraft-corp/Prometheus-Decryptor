win32:
	GOOS=windows GOARCH=386 go build -tags windows -o prometheus_decrypt_32.exe

win64:
	GOOS=windows GOARCH=amd64 go build -tags windows -o prometheus_decrypt_64.exe

linux:
	GOOS=linux GOARCH=amd64 go build

win32GUI:
	set GOARCH=386
	go build -tags windows,gui -ldflags="-H windowsgui -w -s" -o prometheus_decrypt_32_GUI.exe

win64GUI:
	set GOARCH=amd64
	go build -tags windows,gui -ldflags="-H windowsgui -w -s" -o prometheus_decrypt_64_GUI.exe

cross_linux: win32 win64 linux

gui: win32GUI win64GUI

all: win32 win64 linux win32GUI win64GUI
