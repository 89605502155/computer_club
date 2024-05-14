build:
	docker build -t computer-club:v1 .
run:
	docker run -it --rm computer-club:v1
run_without_docker:
	go run cmd/main.go