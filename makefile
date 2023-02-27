start_server:
	go run main.go
docker:
	cd compose && docker-compose up -d
test:
	go test -v tests/*.go
pull:
	git pull origin master
deploy: pull start_server