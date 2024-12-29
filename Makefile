.PHONY: images network container remove build clean run-tests

build:
	go build -o cmd/broker/broker cmd/broker/main.go
	go build -o cmd/auth/auth cmd/auth/main.go
	go build -o cmd/file/file cmd/file/main.go

network:
	docker network create --driver bridge --subnet 10.0.1.0/24 dmz
	docker network create --driver bridge --subnet 10.0.2.0/24 srv
	docker network create --driver bridge --subnet 10.0.3.0/24 dev

images: build
	docker build --rm -f Dockerfile --tag miubuntu .
	docker build --rm -f cmd/broker/Dockerfile --tag miubuntu-broker cmd/broker
	docker build --rm -f cmd/auth/Dockerfile --tag miubuntu-auth cmd/auth
	docker build --rm -f cmd/file/Dockerfile --tag miubuntu-file cmd/file
	docker build --rm -f router/Dockerfile --tag miubuntu-router router
	docker build --rm -f work/Dockerfile --tag miubuntu-work work
	docker build --rm -f jump/Dockerfile --tag miubuntu-jump jump

remove:
	docker stop broker auth file router work jump
	docker network rm dmz
	docker network rm srv
	docker network rm dev

clean:
	rm -f cmd/broker/broker
	rm -f cmd/auth/auth
	rm -f cmd/file/file
	docker rmi miubuntu miubuntu-broker miubuntu-auth miubuntu-file miubuntu-router miubuntu-work miubuntu-jump

container: network 
	docker run --privileged --rm  -ti -d --ip 10.0.1.4 --network dmz --dns 8.8.8.8 --name broker --hostname broker miubuntu-broker

	docker run --privileged --rm -ti -d \
		--name auth --dns 8.8.8.8 --hostname auth --ip 10.0.2.3 --network srv miubuntu-auth

	docker run --privileged --rm -ti -d \
		--name file --hostname file --dns 8.8.8.8 --ip 10.0.2.4 --network srv \
		miubuntu-file
	
	docker run --privileged --rm -ti -d --name work --dns 8.8.8.8 --hostname work --network dev --ip 10.0.3.3 miubuntu-work

	docker run --privileged --rm -ti -d --name jump --hostname jump --dns 8.8.8.8 --network dmz --ip 10.0.1.3 miubuntu-jump
	
	docker run --privileged --rm -ti -d --name router -p 5000:5000 --hostname router --dns 8.8.8.8 --dns 1.1.1.1  miubuntu-router
	docker network connect dmz router
	docker network connect srv router
	docker network connect dev router

run-tests:
	pytest test_https.py