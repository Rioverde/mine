.PHONY: run log restart stop clean psql

run:
	docker compose up --build -d

log:
	docker compose logs -f mine

restart:
	docker compose restart mine

stop:
	docker compose down

clean:
	docker compose down -v

psql:
	docker compose exec storage psql -U mine -d mine
