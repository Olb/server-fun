ifneq (,$(wildcard .env))
	include .env
	export
endif

run-sql:
	docker-compose --env-file .env up --build postgres api

run-mongo:
	docker-compose --env-file .env up -d mongo
	sleep 5
	docker exec -i $(MONGO_CONTAINER) mongosh --quiet blog --eval 'db.posts.insertMany([{ "_id": 1, "title": "Post 1", "body": "Content 1", "created_at": new Date(), "updated_at": new Date() }, { "_id": 2, "title": "Post 2", "body": "Content 2", "created_at": new Date(), "updated_at": new Date() }])'
	DB_TYPE=mongodb MONGODB_URL="$(MONGODB_URL)" go run cmd/main.go

stop:
	docker-compose down

clean:
	docker-compose down -v
