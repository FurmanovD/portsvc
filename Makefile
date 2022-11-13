SHELL := /bin/sh

export APP_NAME=portsvc
export DOCKER_BUILD=docker/portsvc/Dockerfile.build
export TEST_CONTAINER=portsvc-tests
# constructed
export APP_VERSION=$(shell cat ./VERSION)
export IMAGE_NAME ?= ${APP_NAME}
export APPBUILDER_IMAGE ?= ${IMAGE_NAME}-builder
export IMAGE_TAG = ${IMAGE_NAME}:${APP_VERSION}
export IMAGE_TAG_LATEST = ${IMAGE_NAME}:latest
# construct a build information to set it to a binary file variables at 'go build' step:
export APP_VERSION=$(shell cat ./VERSION)
export BUILD_TIME=$$(date -u "+%F_%T")
export GIT_COMMIT=$$(git log -1 --format="%H")

generate:
	go generate ./...

appbuilder-build: 
	# the base image will be built the first time only
	@echo "Check and build if not exist a base App Builder image"
	@docker image inspect ${APPBUILDER_IMAGE} > /dev/null || docker build \
		--tag ${APPBUILDER_IMAGE} \
		--file docker/portsvc/Dockerfile.appbuilder .

appbuilder-clean: 
	@docker rmi ${APPBUILDER_IMAGE}

appbuilder-rebuild: appbuilder-clean appbuilder-build

image-build: appbuilder-build
# the application image will be built one time only
	@echo "APP_VERSION=${APP_VERSION}"
	@echo "BUILD_TIME=${BUILD_TIME}"
	@echo "GIT_COMMIT=${GIT_COMMIT}"
	@echo "check and build if not exist a ${IMAGE_TAG} image"
	@docker image inspect ${IMAGE_TAG} > /dev/null || docker build --no-cache \
		--build-arg APPBUILDER_IMAGE=${APPBUILDER_IMAGE} \
		--build-arg APP_NAME=${APP_NAME} \
		--build-arg APP_VERSION=${APP_VERSION} \
		--build-arg BUILD_TIME=${BUILD_TIME} \
		--build-arg GIT_COMMIT=${GIT_COMMIT} \
		--tag ${IMAGE_TAG} \
		--tag ${IMAGE_TAG_LATEST} \
		--file ${DOCKER_BUILD} .

image-rebuild: image-clean image-build

image-clean:
	@docker rm ${APP_NAME} || echo "container ${APP_NAME} Ok"
	@docker rmi -f ${IMAGE_TAG} || echo "image ${IMAGE_TAG} Ok"
	@docker rmi -f ${IMAGE_TAG_LATEST} || echo "image ${IMAGE_TAG_LATEST} Ok"

image-push: 
	@docker push ${IMAGE_TAG}
	@docker push ${IMAGE_TAG_LATEST}

# builds a binary file
build: clean image-build
	@echo "building ${APP_NAME} (commit:${GIT_COMMIT})"
	@CID=$$(docker create ${IMAGE_TAG}) && \
	docker cp $${CID}:/app/${APP_NAME} ${APP_NAME} && \
	docker rm $${CID}

clean: stop image-clean
	@rm -f ${APP_NAME} || echo ""
	@echo "clean complete"

run: stop image-rebuild run-prebuilt

run-prebuilt:
	@echo "make sure you have a recently built ${IMAGE_TAG} image or run 'make image-rebuild' if not sure"
	@docker image inspect ${IMAGE_TAG} > /dev/null || docker pull ${IMAGE_TAG} || echo ">>> ${IMAGE_TAG} image has not been found! login to a docker repository or run 'make image-build' first"
	IMAGE_TAG=${IMAGE_TAG} docker-compose -f docker/docker-compose.local.yaml up -d

run-deps:
	IMAGE_TAG=${IMAGE_TAG} docker-compose -f docker/docker-compose.local.yaml up -d mysqldb
	
test-docker: 
	docker-compose -f docker/docker-compose.local.yaml \
	 -f docker/docker-compose.tests.yaml up -d tests \
	&& docker logs -f ${TEST_CONTAINER} \
	&& docker-compose -f docker/docker-compose.local.yaml \
	  -f docker/docker-compose.tests.yaml down

stop:
	docker-compose -f docker/docker-compose.local.yaml down --remove-orphans

test: 
	@go test ./... -p=1 -v -cover -failfast

models:
	sqlboiler -c sqlboiler.toml -d -o ./internal/pkg/db/automodel/ mysql

lint:
	golangci-lint run ./...

.PHONY: appbuilder-rebuild image-build image-rebuild image-clean image-push build clean run run-prebuilt run-deps test-docker stop test models generate lint 

