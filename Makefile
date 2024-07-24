jwt-key:
	ssh-keygen -t rsa -b 4096 -m PEM -f ./config/banking_jwtRS256.key

jwt-pub:
	openssl rsa -in ./config/banking_jwtRS256.key -pubout -outform PEM -out ./config/banking_jwtRS256.key.pub

jwt: jwt-key jwt-pub

env: 
	cp ./config/.env.example ./config/.env

temp:
	mkdir temp && mkdir temp/db && mkdir temp/kv && mkdir temp/grafana

compose:
	docker-compose -f ./config/docker-compose.yml up -d

compose-build:
	docker-compose -f ./config/docker-compose.yml up -d --build --remove-orphans

compose-down:
	docker-compose -f ./config/docker-compose.yml down

build-app:
	cd apps/banking && docker build -t github.com/9ssi7/banking:latest .

secret-register:
	docker secret create banking_private_key ./config/banking_jwtRS256.key
	docker secret create banking_public_key ./config/banking_jwtRS256.key.pub

network:
	docker network create --driver overlay --attachable banking

start-app:
	docker service create --name 9ssi7banking --publish 4000:4000 --secret banking_private_key --secret banking_public_key --replicas 3 --env-file ./config/.env --network banking github.com/9ssi7/banking:latest

stop-app:
	docker service rm 9ssi7banking

clean:
	rm -rf temp
	rm -rf config/banking_jwtRS256.key
	rm -rf config/banking_jwtRS256.key.pub

clean-docker:
	docker service rm 9ssi7banking
	docker secret rm banking_private_key
	docker secret rm banking_public_key
	docker network rm banking
	docker rmi github.com/9ssi7/banking:latest

clean-all: clean clean-docker

reqs: temp jwt-key jwt-pub env network compose secret-register

burn: temp jwt-key jwt-pub env network compose secret-register build-app start-app

stop: compose-down stop-app clean

reload: stop-app build-app start-app

.PHONY: jwt-key jwt-pub jwt temp compose compose-down build-app secret-register network start-app stop-app clean start stop burn