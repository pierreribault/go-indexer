dev:
	docker-compose up -d --force-recreate
	docker-compose logs -f go-indexer

prod:
	docker-compose down
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
	docker-compose logs -f go-indexer