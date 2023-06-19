.PHONY: build maria

build:
	go build \
		-ldflags "-X main.buildcommit=`git rev-parse --short HEAD` \
		-X main.buildtime=`date -u '+%Y-%m-%d_%I:%M:%S%p'`" \
		-o main



maria: 
	docker run -p 127.0.0.1:3306:3306 --name some-mariadb \
	-e MARIADB_ROOT_PASSWORD=passwordja -e MARIADB_DATABASE=mydb -d mariadb:latest
