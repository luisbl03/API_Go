.PHONY: build network container remove 

network:
	docker network create --driver bridge --subnet 10.0.1.0/24 dmz
	docker network create --driver bridge --subnet 10.0.2.0/24 srv
	docker network create --driver bridge --subnet 10.0.3.0/24 dev

build:
	docker build --rm -f Dockerfile --tag miubuntu .
	docker build --rm -f cmd/broker/Dockerfile --tag miubuntu-broker cmd/broker
	docker build --rm -f cmd/auth/Dockerfile --tag miubuntu-auth cmd/auth
	docker build --rm -f cmd/file/Dockerfile --tag miubuntu-file cmd/file
	docker build --rm -f router/Dockerfile --tag miubuntu-router router
	docker build --rm -f work/Dockerfile --tag miubuntu-work work
	docker build --rm -f jump/Dockerfile --tag miubuntu-jump jump

remove:
	docker stop broker auth file router work jump
	docker rmi miubuntu
	docker rmi miubuntu-broker
	docker rmi miubuntu-auth
	docker rmi miubuntu-file
	docker rmi miubuntu-work
	docker rmi miubuntu-router
	docker rmi miubuntu-jump
	docker network rm dmz
	docker network rm srv
	docker network rm dev

container: network build
	docker run --privileged --rm -ti -d --ip 10.0.1.4 --network dmz --name broker --hostname broker miubuntu-broker

	docker run --privileged --rm -ti -d \
		--name auth --hostname auth --ip 10.0.2.3 --network srv miubuntu-auth

	docker run --privileged --rm -ti -d \
		--name file --hostname file --ip 10.0.2.4 --network srv \
		miubuntu-file
	
	docker run --privileged --rm -ti -d --name work --hostname work --network dev --ip 10.0.3.3 miubuntu-work

	docker run --privileged --rm -ti -d --name jump --hostname jump --network dmz --ip 10.0.1.3 miubuntu-jump
	
	docker run --privileged --rm -ti -d --name router --hostname router miubuntu-router
	docker network connect dmz router
	docker network connect srv router
	docker network connect dev router