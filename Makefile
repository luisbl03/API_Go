.PHONY: build network container remove clean

network:
	docker network create --driver bridge --subnet 10.0.1.0/24 dmz
	docker network create --driver bridge --subnet 10.0.2.0/24 srv
	docker network create --driver bridge --subnet 10.0.3.0/24 dev

build:
	docker build --rm -f Dockerfile --tag miubuntu .
	docker build --rm -f cmd/broker/Dockerfile --tag miubuntu-broker cmd/broker
	docker build --rm -f cmd/auth/Dockerfile --tag miubuntu-auth cmd/auth
	docker build --rm -f cmd/file/Dockerfile --tag miubuntu-file cmd/file