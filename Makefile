include .env
export

MIGRATIONS_DIR=./db/migrations

run:
	go run ./cmd/api/main.go

db-up:
	docker compose --env-file .env -f deployments/compose.yaml up -d

db-down:
	docker compose --env-file .env -f deployments/compose.yaml down

migrations-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up;

migrations-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

migrations-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" status
	
new_mg:
	@name="$(filter-out $@,$(MAKECMDGOALS))"; \
	last=$$(ls $(MIGRATIONS_DIR)/*.sql 2>/dev/null | grep -E '/[0-9]{5}_' | sort | tail -n 1); \
	if [ -z "$$last" ]; then \
		next="00001"; \
	else \
		last_base=$$(basename "$$last"); \
		next=$$(printf "%05d" $$((10#$${last_base:0:5} + 1))); \
	fi; \
	goose -dir $(MIGRATIONS_DIR) create "$$name" sql; \
	file=$$(ls -t $(MIGRATIONS_DIR)/*.sql | head -n 1); \
	base=$$(basename "$$file"); \
	mv "$$file" "$(MIGRATIONS_DIR)/$${next}_$${base}"; \
	echo "Created: $(MIGRATIONS_DIR)/$${next}_$${base}"

%:
	@:

sqlc:
	sqlc generate
