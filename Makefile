VERSION=development

default:
	@echo "=============Building Windows GUI binary============="
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags desktop,production -ldflags="-H windowsgui -s -w -X 'main.Version=$(VERSION)'" -o dist/simple-proxy-gui.exe .

build-all:
	@echo "=============Building binaries============="

	# Windows 386
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -tags desktop,production -ldflags="-H windowsgui -s -w -X 'main.Version=$(VERSION)'" -o dist/windows_386/simple-proxy-gui.exe .
	cp LICENSE dist/windows_386/LICENSE

	# Windows amd64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags desktop,production -ldflags="-H windowsgui -s -w -X 'main.Version=$(VERSION)'" -o dist/windows_amd64/simple-proxy-gui.exe .
	cp LICENSE dist/windows_amd64/LICENSE

format:
	@echo "=============Formatting============="
	gofmt -s -w .
	go mod tidy

test:
	@echo "=============Running unit tests============="
	go test -v ./...
