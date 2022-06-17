BUILD_DATE := $(shell date -u '+%Y%m%d')
COMMIT := $(shell git describe --tags --always --dirty)
IMAGE_VERSION := $(BUILD_DATE)-$(COMMIT)

BOT_IMAGE_NAME := r.internal.jaxsax.co/tapeworm/botv2:$(IMAGE_VERSION)
REMIX_IMAGE_NAME := r.internal.jaxsax.co/tapeworm/botv2_remix:$(IMAGE_VERSION)

deploy-bot: image-build image-push
	echo 'Hi'

deploy-remix-fe: image-build-remix image-push-remix

image-push-remix:
	docker push $(REMIX_IMAGE_NAME)

image-build-remix:
	cd app-remix && DOCKER_BUILDKIT=1 docker build -t $(REMIX_IMAGE_NAME) .

image-build:
	DOCKER_BUILDKIT=1 docker build -t $(BOT_IMAGE_NAME) .

image-exec:
	docker run -it --rm $(BOT_IMAGE_NAME) bash	

image-push:
	docker push $(BOT_IMAGE_NAME)

dev:
	go run ./cmd/simplifiedbot
