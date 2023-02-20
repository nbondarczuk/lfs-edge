DOCKER_IMAGE=lfs-edge-proxy

# Create a docker image for the service.
docker-image:
	docker build -t $(DOCKER_IMAGE) -f Dockerfile ../..

clean:
	@docker rmi -f $(DOCKER_IMAGE)

.PHONY: docker-image clean
.SILENT:
