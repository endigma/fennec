build:
	@go build -o bin/fennec

portable:
	@CGO_ENABLED=0 go build -o bin/fennec

docker:
	@docker build -t endigma/fennec:latest .

compose:
	@docker-compose build
	@docker-compose up