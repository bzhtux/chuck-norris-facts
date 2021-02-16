APP_NAME = "lcc"
BIN_DIR = "${HOME}/bin"
TMP_DIR := $(shell mktemp -d)

.PHONY: all clean remove help

all: ## ⚆ build and install Go app
all: post-install



get-last-tag: ## ⚆ Get last git tag
	@git tag -l |sort -d -n -r |head -n 1

get-help: ## ⚆ Get last git tag
	@echo "make G_TAG='NEW_GIT_TAG' G_MSG='GIT_TAG_MSG'"

tag: ## ⚆ Git tag current commit with variable G_TAG=${NEW_TAG} G_MSG=${G_MSG}
	@echo "TAG: ${G_TAG}"
	@echo "MSG: ${G_MSG}"
	#git tag -a ${G_TAG} -m "${G_MSG} "

push: ## ⚆ Git push with tags
	@echo "Git push with tags"

post-install: install
	@echo "✨  Cleaning working dir"
	@cp ${TMP_DIR}/${MAC_DIR}/${APP_NAME} ./build/${MAC_DIR}
	@cp ${TMP_DIR}/${LINUX_DIR}/${APP_NAME} ./build/${LINUX_DIR}
	@rm -rf ${TMP_DIR}

pre-build:
	@echo "🛠️   Start pre-build"
	@echo "✨  Deleting ${APP_NAME} from ${BIN_DIR}"
	@if [ -d "./build" ];then rm -rf ./build ;fi && mkdir -p ./build/{${MAC_DIR},${LINUX_DIR},${WINDOWS_DIR}}
	@if [ -f "${BIN_DIR}/${APP_NAME}" ];then rm -f "${BIN_DIR}/${APP_NAME}";fi
	@echo "📂  Creating TMP_DIR"
	@if [ ! -d "${TMP_DIR}" ];then echo "⚙️ mktemp" && mkdir -p ${TMP_DIR} ;fi
	@mkdir -p ${TMP_DIR}/{${MAC_DIR},${LINUX_DIR},${WINDOWS_DIR}}

build: pre-build
	@echo "📦  Building app ${APP_NAME} for 🍏"
	@GOOS=darwin GOARCH=amd64 go build -o ${TMP_DIR}/${MAC_DIR}/${APP_NAME} .
	@echo "📦  Building app ${APP_NAME} for 🐧"
	@GOOS=linux GOARCH=386 go build -o ${TMP_DIR}/${LINUX_DIR}/${APP_NAME} .

install: build
	@echo "🚀  Installing app ${APP_NAME} into ${BIN_DIR}"
	@cp ${TMP_DIR}/${MAC_DIR}/${APP_NAME} "${BIN_DIR}/${APP_NAME}"
	@echo "🥳   Enjoy lcc"

clean: ## ⚆ clean binary from all locations
	@echo "✨   Cleaning working dir"
	@rm -f ${APP_NAME}
	@rm -rf ${TMP_DIR}

# .PHONY: remove
remove: ## ⚆ Remove ${APP_NAME} from ${BIN_DIR}
	@echo "✨  Deleting ${APP_NAME} from ${BIN_DIR}"
	@rm -f "${BIN_DIR}/${APP_NAME}"

.PHONY: help
help: ## ⚆ display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.DEFAULT_GOAL := help

