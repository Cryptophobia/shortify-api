ifndef TAG
$(error The TAG variable is missing.)
endif

ifndef ENV
$(error The ENV variable is missing.)
endif

ifeq ($(filter $(ENV),test dev stag prod),)
$(error The ENV variable is invalid.)
endif

ifeq (,$(filter $(ENV),test dev))
COMPOSE_FILE_PATH := -f docker-compose.yml
endif

IMAGE := cryptophobia/shortify-api

build:
	$(info Make: Building "$(ENV)" environment images.)
	@TAG=$(TAG) docker-compose build --no-cache
	@make -s clean

start:
	$(info Make: Starting "$(ENV)" environment containers.)
	@TAG=$(TAG) docker-compose ${COMPOSE_FILE_PATH} up -d

stop:
	$(info Make: Stopping "$(ENV)" environment containers.)
	@docker-compose stop

restart:
	$(info Make: Restarting "$(ENV)" environment containers.)
	@make -s stop
	@make -s start

push:
	$(info Make: Pushing "$(TAG)" tagged image.)
	@docker push $(IMAGE):$(TAG)

pull:
	$(info Make: Pulling "$(TAG)" tagged image.)
	@docker pull $(IMAGE):$(TAG)

clean:
	@docker system prune --volumes --force
