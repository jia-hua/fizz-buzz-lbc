
build:
	docker build -t fizzbuzz -f ./Dockerfile-production .

run: 
	docker-compose up fizzbuzz

runAll: 
	docker-compose up