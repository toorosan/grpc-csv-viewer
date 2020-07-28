linters:
	docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.28.0 golangci-lint run -v

rebuild:
	cd misc && docker-compose -f docker-compose-everything.yml build --force-rm

run:
	cd misc && docker-compose -f docker-compose-everything.yml up

run-dev:
	cd misc && docker-compose -f docker-compose-everything-dev.yml up
