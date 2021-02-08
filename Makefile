binary: go-tidy go-test go-package go-build

go-package:
	@echo "  >  Packaging static files..."
	pkger

go-build:
	@echo "  >  Building binary..."
	go build -o probr cmd/main.go
	@echo "  >  Removing pkged.go to avoid accidental re-use of old files..."
	rm pkged.go

go-test:
	@echo "  >  Validating code..."
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