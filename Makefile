win32:
	GOOS=windows GOARCH=386 go build -tags windows -o thanos_decrypt_32.exe

win64:
	GOOS=windows GOARCH=amd64 go build -tags windows -o thanos_decrypt_64.exe

linux:
	go build

win32GUI:
	GOOS=windows GOARCH=386 go build -tags windows,gui -ldflags="-H windowsgui -w -s" -o thanos_decrypt_32.exe

win64GUI:
	GOOS=windows GOARCH=amd64 go build -tags windows,gui -ldflags="-H windowsgui -w -s" -o thanos_decrypt_64.exe

all: win32 win64 linux win32GUI win64GUI

