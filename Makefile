win32:
	GOOS=windows GOARCH=386 go build -tags windows -o thanos_decrypt_32.exe thanos_decrypt.go

win64:
	GOOS=windows GOARCH=amd64 go build -tags windows -o thanos_decrypt_64.exe thanos_decrypt.go

linux:
	go build thanos_decrypt.go

all: win32 win64 linux

