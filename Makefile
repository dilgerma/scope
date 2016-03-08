.PHONY: all deps static clean client-lint client-test client-sync backend frontend shell lint

# If you can use Docker without being root, you can `make SUDO= <target>`
SUDO=$(shell (echo "$$DOCKER_HOST" | grep "tcp://" >/dev/null) || echo "sudo -E")
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
CODECGEN_DIR=vendor/github.com/ugorji/go/codec/codecgen
CODECGEN_EXE=$(CODECGEN_DIR)/codecgen
GET_CODECGEN_DEPS=$(shell find $(1) -maxdepth 1 -type f -name '*.go' -not -name '*_test.go' -not -name '*.codecgen.go' -not -name '*.generated.go')
CODECGEN_TARGETS=report/report.codecgen.go render/render.codecgen.go render/detailed/detailed.codecgen.go
RM=--rm
RUN_FLAGS=-ti
BUILD_IN_CONTAINER=true
GO ?= env GO15VENDOREXPERIMENT=1 GOGC=off go
GO_BUILD_INSTALL_DEPS=-i
GO_BUILD_TAGS='netgo unsafe'
GO_BUILD_FLAGS=$(GO_BUILD_INSTALL_DEPS) -ldflags "-extldflags \"-static\" -X main.version=$(SCOPE_VERSION)" -tags $(GO_BUILD_TAGS)

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

$(SCOPE_EXE): $(shell find ./ -path ./vendor -prune -o -type f -name *.go) prog/static.go $(CODECGEN_TARGETS)

report/report.codecgen.go: $(call GET_CODECGEN_DEPS,report/) $(CODECGEN_EXE)
render/render.codecgen.go: $(call GET_CODECGEN_DEPS,render/) $(CODECGEN_EXE)
render/detailed/detailed.codecgen.go: $(call GET_CODECGEN_DEPS,render/detailed/) $(CODECGEN_EXE)
$(CODECGEN_EXE): $(CODECGEN_DIR)/*.go
static: prog/static.go
prog/static.go: client/build/app.js

ifeq ($(BUILD_IN_CONTAINER),true)

$(SCOPE_EXE) $(RUNSVINIT) $(CODECGEN_TARGETS) $(CODECGEN_EXE) lint tests shell prog/static.go: $(SCOPE_BACKEND_BUILD_UPTODATE)
	@mkdir -p $(shell pwd)/.pkg
	$(SUDO) docker run $(RM) $(RUN_FLAGS) \
		-v $(shell pwd):/go/src/github.com/dilgerma/scope \
		-v $(shell pwd)/.pkg:/go/pkg \
		--net=host \
		-e GOARCH -e GOOS -e CIRCLECI -e CIRCLE_BUILD_NUM -e CIRCLE_NODE_TOTAL \
		-e CIRCLE_NODE_INDEX -e COVERDIR -e SLOW \
		$(SCOPE_BACKEND_BUILD_IMAGE) SCOPE_VERSION=$(SCOPE_VERSION) GO_BUILD_INSTALL_DEPS=$(GO_BUILD_INSTALL_DEPS) $@

else

$(SCOPE_EXE): $(SCOPE_BACKEND_BUILD_UPTODATE)
	time $(GO) build $(GO_BUILD_FLAGS) -o $@ ./$(@D)
	@strings $@ | grep cgo_stub\\\.go >/dev/null || { \
	        rm $@; \
	        echo "\nYour go standard library was built without the 'netgo' build tag."; \
	        echo "To fix that, run"; \
	        echo "    sudo go clean -i net"; \
	        echo "    sudo go install -tags netgo std"; \
	        false; \
	    }

%.codecgen.go: $(SCOPE_BACKEND_BUILD_UPTODATE)
	cd $(@D) && env -u GOARCH -u GOOS GOGC=off $(shell pwd)/$(CODECGEN_EXE) -rt $(GO_BUILD_TAGS) -u -o $(@F) $(notdir $(call GET_CODECGEN_DEPS,$(@D)))

$(CODECGEN_EXE): $(SCOPE_BACKEND_BUILD_UPTODATE)
	env -u GOARCH -u GOOS $(GO) build -tags $(GO_BUILD_TAGS) -o $@ ./$(@D)

$(RUNSVINIT): $(SCOPE_BACKEND_BUILD_UPTODATE)
	time $(GO) build $(GO_BUILD_FLAGS) -o $@ ./$(@D)

shell: $(SCOPE_BACKEND_BUILD_UPTODATE)
	/bin/bash

tests: $(SCOPE_BACKEND_BUILD_UPTODATE)
	./tools/test -no-go-get

lint: $(SCOPE_BACKEND_BUILD_UPTODATE)
	./tools/lint .

prog/static.go: $(SCOPE_BACKEND_BUILD_UPTODATE)
	esc -o $@ -prefix client/build client/build

endif

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
	$(GO) clean ./...
	# Don't actually rmi the build images - rm'ing the .uptodate files is enough to ensure
	# we rebuild the images, and rmi'ing the images causes us to have to redownload a lot of stuff.
	# $(SUDO) docker rmi $(SCOPE_UI_BUILD_IMAGE) $(SCOPE_BACKEND_BUILD_IMAGE) >/dev/null 2>&1 || true
	rm -rf $(SCOPE_EXPORT) $(SCOPE_UI_BUILD_UPTODATE) $(SCOPE_BACKEND_BUILD_UPTODATE) \
		$(SCOPE_EXE) $(RUNSVINIT) prog/static.go client/build/app.js docker/weave .pkg \
		$(CODECGEN_TARGETS) $(CODECGEN_EXE)

deps:
	$(GO) get -u -f -tags $(GO_BUILD_TAGS) \
		github.com/FiloSottile/gvt \
		github.com/mattn/goveralls \
		github.com/weaveworks/github-release
