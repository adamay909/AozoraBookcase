default: 
	go build

linux:
	GOOS=linux GOARCH=amd64 go build -o aozoraBookcase_LINUXamd64

macintel:
	GOOS=darwin GOARCH=amd64 go build -o aozoraBookcase_INTEL

applesilicon:
	GOOS=darwin GOARCH=arm64 go build -o aozaraBookcase_APPLESILICON

windows:
	GOOS=windows GOARCH=amd64 go build -o aozoraBookcase.exe

all: linux macintel applesilicon windows
