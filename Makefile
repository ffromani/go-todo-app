all: binaries

outdir: 
	@mkdir -p _out

clean:
	@rm -rf _out

binaries: outdir
	go build -v -o _out/todo cmd/main.go
