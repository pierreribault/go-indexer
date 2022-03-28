dev:
	docker-compose up -d --force-recreate

prod:
	docker-compose down
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build