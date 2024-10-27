all: binaries

outdir: 
	@mkdir -p _out

clean:
	@rm -rf _out coverage.out

binaries: outdir
	go build -v -o _out/todo cmd/main.go

test-unit:
	go test -coverprofile=coverage.out ./...

coverage.out: test-unit

cover-view: coverage.out
	go tool cover -html=coverage.out

test-e2e:
	ginkgo -v ./e2e/...
