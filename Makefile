# Load environment variables from .env
include .env
export $(shell sed 's/=.*//' .env)
# Container name and image
CONTAINER_NAME=postgres-container
IMAGE_NAME=postgres:15
VOLUME_NAME=postgres_data
# Start PostgreSQL in Docker
start-db:
	@echo "Starting PostgreSQL container..."
	@docker run --name $(CONTAINER_NAME) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p 5432:5432 \
		-v $(VOLUME_NAME):/var/lib/postgresql/data \
		-d $(IMAGE_NAME)
	@echo "PostgreSQL container started."
# Stop the PostgreSQL container
stop-db:
	@echo "Stopping PostgreSQL container..."
	@docker stop $(CONTAINER_NAME)
	@echo "PostgreSQL container stopped."
# Create the database
create-db:
	@echo "Creating database $(POSTGRES_DB)..."
	@PGPASSWORD=$(POSTGRES_PASSWORD) psql -U $(POSTGRES_USER) -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -c "CREATE DATABASE $(POSTGRES_DB);"
	@echo "Database $(POSTGRES_DB) created successfully."
# Create the table in the database
create-table:
	@echo "Creating table 'messages' in $(POSTGRES_DB)..."
	@PGPASSWORD=$(POSTGRES_PASSWORD) psql -U $(POSTGRES_USER) -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -d $(POSTGRES_DB) -c "\
	CREATE TABLE IF NOT EXISTS messages ( \
		message_id SERIAL PRIMARY KEY, \
		timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, \
		message TEXT NOT NULL, \
		user_id INT NOT NULL, \
		sender_id INT NOT NULL, \
		status VARCHAR(10) NOT NULL DEFAULT 'unread' CHECK (status IN ('unread', 'read')) \
	);"
	@echo "Table 'messages' created successfully."
# Clean the database by dropping the messages table
drop-table:
	@echo "Dropping table 'messages' from $(POSTGRES_DB)..."
	@PGPASSWORD=$(POSTGRES_PASSWORD) psql -U $(POSTGRES_USER) -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -d $(POSTGRES_DB) -c "DROP TABLE IF EXISTS messages;"
	@echo "Table 'messages' dropped successfully."
# Clean up and drop the database
drop-db:
	@echo "Dropping database $(POSTGRES_DB)..."
	@PGPASSWORD=$(POSTGRES_PASSWORD) psql -U $(POSTGRES_USER) -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -c "DROP DATABASE IF EXISTS $(POSTGRES_DB);"
	@echo "Database $(POSTGRES_DB) dropped successfully."
# Run both create-db and create-table in one go
setup: create-db create-table
.PHONY: create-db create-table drop-table drop-db setup