# we assume all test machines have either docker-compose or
# the compose plugin in docker command.
DOCKER_COMPOSE=docker-compose
HAS_DOCKER_COMPOSE := $(shell command -v docker-compose 2> /dev/null)
ifndef HAS_DOCKER_COMPOSE
  DOCKER_COMPOSE=docker compose
  DOCKER_COMPOSE_QUIET_PULL=--quiet-pull
endif

include .env

start:
	$(DOCKER_COMPOSE) -p$(PROJECT) up $(DOCKER_COMPOSE_QUIET_PULL) -d

test: stop_test start sync_test proxy_test

fds:
	$(DOCKER_COMPOSE) -p$(PROJECT) \
		-f docker-compose-fds.yml up $(DOCKER_COMPOSE_QUIET_PULL) -d

fds_deps:
	$(DOCKER_COMPOSE) -p$(PROJECT) \
		-f docker-compose-fds.yml \
		-f docker-compose-fds-deps.yml \
		up $(DOCKER_COMPOSE_QUIET_PULL) -d

common_test:
	$(DOCKER_COMPOSE) -p$(PROJECT) \
		-f docker-compose-common-test.yml up $(DOCKER_COMPOSE_QUIET_PULL) \
		--exit-code-from $(COMMON_TEST)

sync_test:
	$(DOCKER_COMPOSE) -p$(PROJECT) \
		-f docker-compose-sync-test.yml up $(DOCKER_COMPOSE_QUIET_PULL) \
		--exit-code-from $(SYNC_TEST)

proxy_test:
	$(DOCKER_COMPOSE) -p$(PROJECT) \
		-f docker-compose-proxy-test.yml up $(DOCKER_COMPOSE_QUIET_PULL) \
		--exit-code-from $(PROXY_TEST)

stop_test:
	-docker rm $(SYNC_TEST).$(DOMAIN)
	-docker rm $(PROXY_TEST).$(DOMAIN)

stop:
	-$(DOCKER_COMPOSE) -p$(PROJECT) \
		-f docker-compose-fds.yml \
		-f docker-compose-fds-dev.yml \
		-f docker-compose-common-test.yml \
		-f docker-compose-sync-test.yml \
		-f docker-compose-proxy-test.yml down

clean: stop
	-docker ps -aqf status=exited | xargs docker rm # remove exited

.SILENT:
.PHONY: start stop test test_cli sync_test proxy_test stop_test clean
