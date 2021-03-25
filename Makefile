binary: go-tidy go-test go-package go-build pkgr-clean
quick: go-package go-build pkgr-clean

go-package:
	@echo "  >  Packaging static files..."
	pkger

go-build:
	@echo "  >  Building binary..."
	go build -o kubernetes cmd/main.go

pkgr-clean:
	@echo "  >  Removing pkged.go to avoid accidental re-use of old files..."
	rm pkged.go

go-test:
	@echo "  >  Validating code..."
	golint ./...
	go vet ./...
	go test ./...

go-tidy:
	@echo "  >  Tidying go.mod ..."
	go mod tidy

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo