BINARY = ustress
GOARCH = amd64

VERSION?=?
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
GITURL=$(shell git config --get remote.origin.url | sed "s/git@//g;s/\.git//g;s/:/\//g")

CURRENT_DIR=$(shell pwd)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})

DOCKER_IMAGE_REPO       ?= peng-docker-prod.metroscales.io/reliability
DOCKER_IMAGE_NAME       ?= ${BINARY}
DOCKER_IMAGE_TAG        ?= latest

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-w -s -X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

# Build the project
all: linux docker


build_docker:
	docker build . -f ./build/build.Dockerfile -t ustress

clean:
	go clean -n
	rm -f ${CURRENT_DIR}/${BINARY}-windows-${GOARCH}.exe
	rm -f ${CURRENT_DIR}/${BINARY}-linux-${GOARCH}

darwin:
	@echo ">> building linux binary"
	CGO_ENABLED=0 GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH} cmd/ustress/main.go ;

linux:
	@echo ">> building linux binary"
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH} cmd/ustress/main.go ;

windows:
	@echo ">> building windows binary"
	GOOS=windows GOARCH=${GOARCH} go build -o ${BINARY}-windows-${GOARCH}.exe cmd/ustress/main.go ;

docker:
	@echo ">> building docker image"
	docker build -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" -f ./build/package/build.Dockerfile .;

docker-run:
	@echo ">> running docker image"
	docker run --rm -p 8080:8080 "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)";

docker-release:
	@echo ">> make realease";
	docker tag $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) peng-docker-prod.metroscales.io/reliability/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG);
	docker push $(DOCKER_IMAGE_REPO)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG);

# release: linux docker

.PHONY: release all linux windows docker
