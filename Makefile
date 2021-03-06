default: build docker

build:
	GOOS=linux go build

docker:
	docker build -t jspc/router .

run: build docker
	docker run -e DOCKER_API_VERSION=1.37 -v /var/run/docker.sock:/var/run/docker.sock -p 0.0.0.0:80:8080 -ti jspc/router
