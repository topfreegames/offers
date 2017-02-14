MY_IP=`ifconfig | grep --color=none -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep --color=none -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1' | head -n 1`

setup:
	@go get github.com/jteeuwen/go-bindata
	@glide install

assets:
	@go-bindata -o migrations/migrations.go -pkg migrations migrations/*.sql

migrate: assets
	@go run main.go migrate -c ./config/local.yaml

drop:
	@-psql -d postgres -h localhost -p 8585 -U postgres -c "SELECT pg_terminate_backend(pid.pid) FROM pg_stat_activity, (SELECT pid FROM pg_stat_activity where pid <> pg_backend_pid()) pid WHERE datname='offers';"
	@psql -d postgres -h localhost -p 8585 -U postgres -f scripts/drop.sql > /dev/null
	@echo "Database created successfully!"

drop-test:
	@-psql -d postgres -h localhost -p 8585 -U postgres -c "SELECT pg_terminate_backend(pid.pid) FROM pg_stat_activity, (SELECT pid FROM pg_stat_activity where pid <> pg_backend_pid()) pid WHERE datname='offers-test';"
	@psql -d postgres -h localhost -p 8585 -U postgres -f scripts/drop-test.sql > /dev/null
	@echo "Test Database created successfully!"

acceptance acc:
	@go test $(go list ./... | grep -v /vendor/) -tags=acceptance

wait-for-pg:
	@until docker exec offers_postgres_1 pg_isready; do echo 'Waiting for Postgres...' && sleep 1; done
	@sleep 2

deps: start-deps wait-for-pg

start-deps:
	@echo "Starting dependencies using HOST IP of ${MY_IP}..."
	@env MY_IP=${MY_IP} docker-compose --project-name offers up -d
	@sleep 10
	@echo "Dependencies started successfully."

stop-deps:
	@env MY_IP=${MY_IP} docker-compose --project-name offers down

test: deps drop-test
	@ginkgo -r -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements .
