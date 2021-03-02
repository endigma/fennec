development:
	@go build -o bin/fennec

docker:
	@docker build -t endigma/fennec:latest .

compose:
	@docker-compose build
	@docker-compose up