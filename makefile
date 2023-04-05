import_database:
	cd data_import && go run import_data.go
start_server:
	go run main.go
docker:
	cd compose && docker-compose up -d
test:
	go test -v tests/*.go
pull:
	git pull origin master
deploy: pull start_server