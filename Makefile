run-everything:
	cd misc && docker-compose -f everything-docker-compose.yml up

run-everything-rebuilt:
	cd misc && docker-compose -f everything-docker-compose.yml up --build --force-recreate
