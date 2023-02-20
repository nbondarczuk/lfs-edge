TEST_RUNNER_IMAGE=docker.github.azc.ext.hp.com/krypton/utils/krypton-test-runner

iso: lint
	./build_iso_image.sh

lint:
	@docker run --rm \
	-v$(CURDIR):/src \
	-w/src \
	$(TEST_RUNNER_IMAGE) \
	shellcheck *.sh

.PHONY: iso lint
