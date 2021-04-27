release: go-release
rc: go-release-candidate
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

go-release-candidate: binary
	@echo "  >  Building release candidate ..."
	go build -o kubernetes -ldflags="-X 'main.GitCommitHash=`git rev-parse --short HEAD`' -X 'main.BuiltAt=`date +%FT%T%z`' -X 'main.Prerelease=rc'" cmd/main.go

go-release: binary
	@echo "  >  Building release ..."
	go build -o kubernetes -ldflags="-X 'main.GitCommitHash=`git rev-parse --short HEAD`' -X 'main.BuiltAt=`date +%FT%T%z`' -X 'main.Prerelease='" cmd/main.go