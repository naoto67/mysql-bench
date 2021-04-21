bench:
	go test ./... -bench . -benchmem -test.run=none -benchtime 3s

docker_up:
	docker-compose up -d --build

docker_down:
	docker-compose down

mysql_console:
	docker-compose exec mysql mysql -uroot bench
