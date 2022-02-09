# Configurable
SHELL=/bin/bash
CI_SUPPORT_REPO="https://$$GITHUB_TOKEN@github.com/RainThief/ci-support.git"
APP_IMAGE_NAME="ghcr.io/rainthief/spa-server"
CI_BRANCH="master"
CI_DIR_NAME="build/${CI_REPO_NAME}"
GIT_URL=github.com/RainThief/spa-server.git
# add to local env if wish to use local ci support folder
LOCAL_CI_PATH?=""

HEALTHCHECK_PORT?=8000

# Dynamic (do not touch)
PROJECT_ROOT="$$(cd "$$(dirname "$${BASH_SOURCE[0]}" )" && pwd)"
CI_REPO_NAME="$$(sed -E  's/(.*)\/([a-z0-9-]+)\.git$$/\2/' <<< ${CI_SUPPORT_REPO})"


build: all


prepare:
	if [ "${LOCAL_CI_PATH}" != "" ]; then \
		make prepare-dev; \
	else \
		if [ ! -d "${CI_DIR_NAME}" ]; then \
			git clone "${CI_SUPPORT_REPO}" --depth 1 --branch ${CI_BRANCH} "${CI_DIR_NAME}"; \
		else \
			pushd "${CI_DIR_NAME}" > /dev/null; \
				git reset --hard; \
				git pull origin ${CI_BRANCH}; \
			popd > /dev/null; \
		fi \
	fi


prepare-dev:
	echo "copying local ci folder into project"
	mkdir -p "${CI_DIR_NAME}"; \
	cp -R ${LOCAL_CI_PATH}/* "$$PWD/${CI_DIR_NAME}"


build-docker:
	make prepare
	export PROJECT_ROOT=${PROJECT_ROOT}; \
	export APP_IMAGE_NAME=${APP_IMAGE_NAME}; \
	./${CI_DIR_NAME}/common/build_image.sh ${args};


all:
	make audit
	make static-analysis
	make unit-test
	make system-test


audit:
	make prepare
	export PROJECT_ROOT=${PROJECT_ROOT}; \
	source ${CI_DIR_NAME}/common/docker.sh; \
	exec_ci_container go 1.16 "${CI_DIR_NAME}/go/audit.sh";


static-analysis:
	make prepare
	export PROJECT_ROOT=${PROJECT_ROOT}; \
	source ${CI_DIR_NAME}/common/docker.sh; \
	exec_ci_container common v1 "${CI_DIR_NAME}/common/static_analysis.sh"; \
	exec_ci_container go 1.16 "${CI_DIR_NAME}/go/static_analysis.sh";


unit-test:
	make prepare
	export PROJECT_ROOT=${PROJECT_ROOT}; \
	source ${CI_DIR_NAME}/common/docker.sh; \
	exec_ci_container go 1.16 "${CI_DIR_NAME}/go/unit_test.sh" "$@"


system-test:
	make build-docker
	export PROJECT_ROOT=${PROJECT_ROOT}; \
	source ${CI_DIR_NAME}/common/docker.sh; \
	echo_warning "system tests not yet implemented"


test:
	make prepare
	make unit-test
	make system-test


run:
	export SCAN="false"; \
	export DOCKER_BUILD_ARGS="--build-arg HEALTH_CHECK_PORT=${HEALTHCHECK_PORT}"; \
	make build-docker
	export PROJECT_ROOT=${PROJECT_ROOT}; \
	source ${CI_DIR_NAME}/common/docker.sh; \
	IMAGE_TAG="$$(docker_branch_tag)"; \
	docker run --rm -it --init -u="$(id -u):$(id -g)" --name "dev-spa-server" \
		--network=host \
		-e HEALTHCHECK_PORT=${HEALTHCHECK_PORT} \
		"${APP_IMAGE_NAME}:$$IMAGE_TAG"


release:
	make prepare
	export CI=${CI}; \
	export GIT_URL=${GIT_URL}; \
	export TAG=true; \
	./${CI_DIR_NAME}/common/release.sh;
	make build-docker
