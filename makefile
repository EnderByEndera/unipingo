import_database:
	cd data_import && go run import_data.go
start_dev_server:
	sudo systemctl stop internet-plus-backend.service
	go run main.go
update_service:
	go build
	sudo systemctl restart internet-plus-backend.service
docker:
	cd compose && docker-compose up -d
test:
	go test -v tests/*.go
pull:
	git pull origin master
deploy: pull start_server