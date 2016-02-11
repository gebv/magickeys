
GOPATH := ${PWD}:${PWD}/vendor
PATH := ${PATH}:${PWD}/vendor/bin

VERSION := "0.1.DEV"
BUILDSTAMP :=`date +%FT%T%z`

GITHASH := `git rev-parse HEAD`

# GOOS:="linux"
# GOARCH:="386"

fmt:
	GOPATH=$(GOPATH) go fmt ./src/...

vet:
	GOPATH=$(GOPATH) go vet ./src/...

build:
	GOPATH=$(GOPATH) PATH=$(PATH) go generate ./src/...
	GOPATH=$(GOPATH) go build -ldflags "-X utils.GitHash=$(GITHASH) -X utils.BuildDate=$(BUILDSTAMP) -X utils.Version=$(VERSION)" -v -o ./bin/app.bin ./src

run:
	GOPATH=$(GOPATH) go run -ldflags "-X utils.GitHash=$(GITHASH) -X utils.BuildDate=$(BUILDSTAMP) -X utils.Version=$(VERSION)" ./src/main.go -stderrthreshold=INFO -v=2

test:
	GOPATH=$(GOPATH) go test ./src/...  -test.bench=. -test.benchmem testing: warning: no tests to run

vendor_clean:
	# find ./src -type d -not -name '*.run' | xargs rm
	rm -Rf ./vendor
	mkdir -p ./vendor
	rm -Rf ./bin
	mkdir -p ./bin
	rm -Rf ./pkg

vendor_update: vendor_get
	rm -rf `find ./vendor -type d -name .git` \
	&& rm -rf `find ./vendor -type d -name .hg` \
	&& rm -rf `find ./vendor -type d -name .bzr` \
	&& rm -rf `find ./vendor -type d -name .svn`

vendor_get: vendor_clean
	GOPATH=${PWD}/vendor go get -u -v \
		github.com/gorilla/mux \
		github.com/golang/glog \
		github.com/lib/pq \
		github.com/gebv/goco \
		gopkg.in/bluesuncorp/validator.v8 \
		github.com/jackc/pgx \
		github.com/satori/go.uuid \