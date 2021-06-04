linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" .
	zip aws-mfa-auth.linux.amd64.zip aws-mfa-auth
	rm -rf aws-mfa-auth
	
darwin:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" .
	zip aws-mfa-auth.darwin.amd64.zip aws-mfa-auth
	rm -rf aws-mfa-auth

windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" .
	zip aws-mfa-auth.windows.amd64.zip aws-mfa-auth.exe
	rm -rf aws-mfa-auth.exe

all: linux darwin windows