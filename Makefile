lint:
	docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.28.0 golangci-lint run -v

rebuild:
	cd misc && docker-compose -f docker-compose-everything.yml build --force-rm

regen-protobuf:
	docker run -v `pwd`:/defs namely/protoc-all -f internal/pkg/csvviewer/csv_viewer.proto -l go -o internal/pkg/csvviewer

run:
	cd misc && docker-compose -f docker-compose-everything.yml up

run-dev:
	cd misc && docker-compose -f docker-compose-everything-dev.yml up
