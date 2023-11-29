build:
	@echo "Exporting for all platorms..."
	env GOOS=darwin GOARCH=arm64 go build -o "./bin/albumcut-darwin-arm"
	env GOOS=darwin GOARCH=amd64 go build -o "./bin/albumcut-darwin-amd"
	env GOOS=windows GOARCH=amd64 go build -o "./bin/albumcut-windows"
	env GOOS=linux GOARCH=arm64 go build -o "./bin/albumcut-linux"


