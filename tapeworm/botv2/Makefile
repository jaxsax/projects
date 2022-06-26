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

create_migration:
	migrate create -ext sql -dir db/migrations $(name)

migrate_up:
	go run ./cmd/migrate --up	
	
migrate_down:
	go run ./cmd/migrate --down

download:
	@echo Download go.mod dependencies
	@go mod download

install-tools: download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %