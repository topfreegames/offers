MY_IP=`ifconfig | grep --color=none -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep --color=none -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1' | head -n 1`
PACKAGES = $(shell glide novendor)
TEST_PACKAGES = $(shell glide novendor | egrep -v features | egrep -v '^[.]$$' | sed 's@\/[.][.][.]@@')

setup:
	@go get github.com/jteeuwen/go-bindata
	@glide install

build:
	@go build $(PACKAGES)
	@go build -o ./bin/khan main.go

assets:
	@go-bindata -o migrations/migrations.go -pkg migrations migrations/*.sql

migrate: assets
	@go run main.go migrate -c ./config/local.yaml

drop:
	@-psql -d postgres -h localhost -p 8585 -U postgres -c "SELECT pg_terminate_backend(pid.pid) FROM pg_stat_activity, (SELECT pid FROM pg_stat_activity where pid <> pg_backend_pid()) pid WHERE datname='offers';"
	@psql -d postgres -h localhost -p 8585 -U postgres -f scripts/drop.sql > /dev/null
	@echo "Database created successfully!"

migrate-test: assets
	@go run main.go migrate -c ./config/test.yaml

drop-test:
	@-psql -d postgres -h localhost -p 8585 -U postgres -c "SELECT pg_terminate_backend(pid.pid) FROM pg_stat_activity, (SELECT pid FROM pg_stat_activity where pid <> pg_backend_pid()) pid WHERE datname='offers-test';"
	@psql -d postgres -h localhost -p 8585 -U postgres -f scripts/drop-test.sql > /dev/null
	@echo "Test Database created successfully!"

acceptance acc:
	@cd features && go test

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

test: deps drop-test unit integration test-coverage-func

unit: clear-coverage-profiles unit-run gather-unit-profiles

integration int: clear-coverage-profiles integration-run gather-integration-profiles

clear-coverage-profiles:
	@find . -name '*.coverprofile' -delete

gather-unit-profiles:
	@mkdir -p '_build'
	@echo "mode: count" > _build/coverage-unit.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/coverage-unit.out; done'

unit-run:
	@ginkgo -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}

gather-integration-profiles:
	@mkdir -p '_build'
	@echo "mode: count" > _build/coverage-integration.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/coverage-integration.out; done'

integration-run:
	@ginkgo -tags integration -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}

merge-profiles:
	@mkdir -p '_build'
	@echo "mode: count" > _build/coverage-all.out
	@bash -c 'for f in $$(find . -name "*.out"); do tail -n +2 $$f >> _build/coverage-all.out; done'

test-coverage-func coverage-func: merge-profiles
	@go tool cover -func=_build/coverage-all.out

test-coverage-html coverage-html: merge-profiles
	@go tool cover -html=_build/coverage-all.out