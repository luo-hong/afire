BASE_PATH=$(shell pwd)
LDFLAGS = "-w -s"

GOOS=$(shell go env GOOS)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

.PHONY: version
.PHONY: manager
.PHONY: resource
.PHONY: clean


manager: version
	@echo "build ${GOOS} afire_manager"
	@echo "build at git "${BRANCH}
ifeq ($(BRANCH),master)
	@cd ${BASE_PATH}/cmd/manager && GOOS=${GOOS} go build -mod=vendor -ldflags ${LDFLAGS} -o ${BASE_PATH}/bin/afire_man
else
	@cd ${BASE_PATH}/cmd/manager && GOOS=${GOOS} go build -race -mod=vendor -o ${BASE_PATH}/bin/afire_man
endif

version:
	@bash ${BASE_PATH}/version/version.sh

resource:
	@go-bindata -o configs/resources.go -pkg=configs configs/*.xml configs/*.json


clean:
	@echo "clean"
	@rm -rf ${BASE_PATH}/bin/*
