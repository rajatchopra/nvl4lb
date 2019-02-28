
.PHONY: build

build: l4controller configurelb

l4controller: cmd/l4controller/_output/bin/l4controller
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o cmd/l4controller/_output/bin/l4controller cmd/l4controller/main.go

configurelb: cmd/configurelb/_output/bin/configurelb
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o cmd/configurelb/_output/bin/configurelb cmd/configurelb/main.go

clean:
	rm -rf cmd/l4controller/_output
	rm -rf cmd/configurelb/_output
