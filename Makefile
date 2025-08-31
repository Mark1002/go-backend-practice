.PHONY: mysql-up mysql-down mysql-logs mysql-shell run build clean

# Docker MySQL commands
mysql-up:
	docker run -d \
		--name mysql-practice \
		-e MYSQL_ROOT_PASSWORD=rootpassword \
		-e MYSQL_DATABASE=practice_db \
		-e MYSQL_USER=appuser \
		-e MYSQL_PASSWORD=apppassword \
		-p 3306:3306 \
		-v mysql-practice-data:/var/lib/mysql \
		mysql:8.0
	@echo "Waiting for MySQL to initialize..."
	@sleep 10
	@echo "Granting super privileges to appuser..."
	@docker exec mysql-practice mysql -u root -prootpassword -e "GRANT ALL PRIVILEGES ON *.* TO 'appuser'@'%' WITH GRANT OPTION; GRANT SUPER ON *.* TO 'appuser'@'%'; GRANT SYSTEM_VARIABLES_ADMIN ON *.* TO 'appuser'@'%'; FLUSH PRIVILEGES;"

mysql-down:
	docker stop mysql-practice && docker rm mysql-practice

mysql-logs:
	docker logs -f mysql-practice

mysql-shell:
	docker exec -it mysql-practice mysql -u root -p

# Migration commands
migrate-up:
	docker exec -i mysql-practice mysql -u appuser -papppassword practice_db < migration/mysql/v0.sql

migrate-down:
	docker exec -i mysql-practice mysql -u appuser -papppassword practice_db -e "DROP TABLE IF EXISTS users; DROP TABLE IF EXISTS orders;"

# Development workflow
dev-setup: mysql-up
	@echo "Waiting for MySQL to be ready..."
	@sleep 10
	@make migrate-up
	@echo "Development environment ready!"

dev-teardown: mysql-down clean

# Volume management
mysql-volume-remove:
	docker volume rm mysql-practice-data