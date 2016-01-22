.PHONY: all deps static clean client-lint client-test client-sync backend frontend

# If you can use Docker without being root, you can `make SUDO= <target>`
SUDO=
DOCKERHUB_USER=dilgerm
SCOPE_EXE=prog/scope
SCOPE_IMAGE=$(DOCKERHUB_USER)/scope
SCOPE_EXPORT=scope.tar
SCOPE_UI_BUILD_IMAGE=$(DOCKERHUB_USER)/scope-ui-build
SCOPE_UI_BUILD_UPTODATE=.scope_ui_build.uptodate
SCOPE_BACKEND_BUILD_IMAGE=$(DOCKERHUB_USER)/scope-backend-build
SCOPE_BACKEND_BUILD_UPTODATE=.scope_backend_build.uptodate
SCOPE_VERSION=$(shell git rev-parse --short HEAD)
DOCKER_VERSION=1.6.2
DOCKER_DISTRIB=docker/docker-$(DOCKER_VERSION).tgz
DOCKER_DISTRIB_URL=https://github.com/dilgerma/weave/blob/master/prog/weaveexec/docker.tgz
RUNSVINIT=vendor/runsvinit/runsvinit
RM=--rm
RUN_FLAGS=-ti
BUILD_IN_CONTAINER=true

all: $(SCOPE_EXPORT)

$(DOCKER_DISTRIB):
	curl -o $(DOCKER_DISTRIB) $(DOCKER_DISTRIB_URL)

docker/weave:
	curl -L https://raw.githubusercontent.com/dilgerma/weave/master/weave -o docker/weave
	chmod u+x docker/weave

$(SCOPE_EXPORT): $(SCOPE_EXE) $(DOCKER_DISTRIB) docker/weave $(RUNSVINIT) docker/Dockerfile docker/run-app docker/run-probe docker/entrypoint.sh
	cp $(SCOPE_EXE) $(RUNSVINIT) docker/
	cp $(DOCKER_DISTRIB) docker/docker.tgz
	$(SUDO) docker build -t $(SCOPE_IMAGE) docker/
	$(SUDO) docker save $(SCOPE_IMAGE):latest > $@

$(RUNSVINIT): vendor/runsvinit/*.go
	go build -o $@ github.com/dilgerma/scope/vendor/runsvinit

$(SCOPE_EXE): $(shell find ./ -type f -name *.go) prog/static.go

ifeq ($(BUILD_IN_CONTAINER),true)
$(SCOPE_EXE) $(RUNSVINIT): $(SCOPE_BACKEND_BUILD_UPTODATE)
	$(SUDO) docker run $(RM) $(RUN_FLAGS) -v $(shell pwd):/gopath1.5/src/github.com/dilgerma/scope -e GOARCH -e GOOS \
		$(SCOPE_BACKEND_BUILD_IMAGE) SCOPE_VERSION=$(SCOPE_VERSION) $@
else
$(SCOPE_EXE): $(SCOPE_BACKEND_BUILD_UPTODATE)
	go build -ldflags "-extldflags \"-static\" -X main.version=$(SCOPE_VERSION)" -tags netgo -o $@ ./$(@D)
	@strings $@ | grep cgo_stub\\\.go >/dev/null || { \
	        rm $@; \
	        echo "\nYour go standard library was built without the 'netgo' build tag."; \
	        echo "To fix that, run"; \
	        echo "    sudo go clean -i net"; \
	        echo "    sudo go install -tags netgo std"; \
	        false; \
	    }

$(RUNSVINIT):
	go build -ldflags "-extldflags \"-static\"" -o $@ ./$(@D)
endif

static: prog/static.go

prog/static.go: client/build/app.js
	esc -o $@ -prefix client/build client/build

ifeq ($(BUILD_IN_CONTAINER),true)
client/build/app.js: $(shell find client/app/scripts -type f) $(SCOPE_UI_BUILD_UPTODATE)
	mkdir -p client/build
	$(SUDO) docker run $(RM) $(RUN_FLAGS) -v $(shell pwd)/client/app:/home/weave/app \
		-v $(shell pwd)/client/build:/home/weave/build \
		$(SCOPE_UI_BUILD_IMAGE) npm run build

client-test: $(shell find client/app/scripts -type f) $(SCOPE_UI_BUILD_UPTODATE)
	$(SUDO) docker run $(RM) $(RUN_FLAGS) -v $(shell pwd)/client/app:/home/weave/app \
		-v $(shell pwd)/client/test:/home/weave/test \
		$(SCOPE_UI_BUILD_IMAGE) npm test

client-lint: $(SCOPE_UI_BUILD_UPTODATE)
	$(SUDO) docker run $(RM) $(RUN_FLAGS) -v $(shell pwd)/client/app:/home/weave/app \
		-v $(shell pwd)/client/test:/home/weave/test \
		$(SCOPE_UI_BUILD_IMAGE) npm run lint

client-start: $(SCOPE_UI_BUILD_UPTODATE)
	$(SUDO) docker run $(RM) $(RUN_FLAGS) --net=host -v $(shell pwd)/client/app:/home/weave/app \
		-v $(shell pwd)/client/build:/home/weave/build \
		$(SCOPE_UI_BUILD_IMAGE) npm start
else
client/build/app.js:
	cd client && npm run build
endif

$(SCOPE_UI_BUILD_UPTODATE): client/Dockerfile client/package.json client/webpack.local.config.js client/webpack.production.config.js client/server.js client/.eslintrc
	$(SUDO) docker build -t $(SCOPE_UI_BUILD_IMAGE) client
	touch $@

$(SCOPE_BACKEND_BUILD_UPTODATE): backend/*
	$(SUDO) docker build -t $(SCOPE_BACKEND_BUILD_IMAGE) backend
	touch $@

clean:
	go clean ./...
	$(SUDO) docker rmi $(SCOPE_UI_BUILD_IMAGE) $(SCOPE_BACKEND_BUILD_IMAGE) >/dev/null 2>&1 || true
	rm -rf $(SCOPE_EXPORT) $(SCOPE_UI_BUILD_UPTODATE) $(SCOPE_BACKEND_BUILD_UPTODATE) \
		$(SCOPE_EXE) $(RUNSVINIT) prog/static.go client/build/app.js docker/weave

ifeq ($(BUILD_IN_CONTAINER),true)
tests: $(SCOPE_BACKEND_BUILD_UPTODATE)
	$(SUDO) docker run $(RM) $(RUN_FLAGS) -v $(shell pwd):/gopath1.5/src/github.com/dilgerma/scope \
		-e GOARCH -e GOOS -e CIRCLECI -e CIRCLE_BUILD_NUM -e CIRCLE_NODE_TOTAL -e CIRCLE_NODE_INDEX -e COVERDIR\
		$(SCOPE_BACKEND_BUILD_IMAGE) tests
else
tests:
	./tools/test -no-go-get
endif

deps:
	go get -u -f -tags netgo \
		github.com/golang/lint/golint \
		github.com/fzipp/gocyclo \
		github.com/mattn/goveralls \
		github.com/mjibson/esc \
		github.com/kisielk/errcheck \
		github.com/weaveworks/github-release
