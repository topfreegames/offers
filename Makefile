MY_IP=`ifconfig | grep --color=none -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep --color=none -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1' | head -n 1`
PACKAGES = $(shell glide novendor)
TEST_PACKAGES = $(shell glide novendor | egrep -v features | egrep -v '^[.]$$' | sed 's@\/[.][.][.]@@')

setup: setup-hooks
	@go get -u github.com/Masterminds/glide/...
	@go get -u github.com/jteeuwen/go-bindata/...
	@go get -u github.com/wadey/gocovmerge
	@glide install

setup-hooks:
	@cd .git/hooks && ln -sf ../../hooks/pre-commit.sh pre-commit

setup-ci:
	@go get github.com/mattn/goveralls
	@go get github.com/onsi/ginkgo/ginkgo
	@${MAKE} setup

build:
	@go build $(PACKAGES)
	@go build -o ./bin/offers main.go

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

test: deps unit integration acceptance test-coverage-func

clear-coverage-profiles:
	@find . -name '*.coverprofile' -delete

unit: drop-test migrate-test clear-coverage-profiles unit-run gather-unit-profiles

unit-run:
	#@LOGXI="*=ERR,dat:sqlx=OFF,dat=OFF" ginkgo -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}
	@ginkgo -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}

gather-unit-profiles:
	@mkdir -p _build
	@echo "mode: count" > _build/coverage-unit.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/coverage-unit.out; done'

integration int: drop-test migrate-test clear-coverage-profiles integration-run gather-integration-profiles

integration-run:
	#@LOGXI="*=ERR,dat:sqlx=OFF,dat=OFF" ginkgo -tags integration -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}
	@ginkgo -tags integration -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}

gather-integration-profiles:
	@mkdir -p _build
	@echo "mode: count" > _build/coverage-integration.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/coverage-integration.out; done'

acceptance acc: drop-test migrate-test clear-coverage-profiles acceptance-run
acceptance-focus acc-focus: drop-test migrate-test clear-coverage-profiles acceptance-run-focus

acceptance-run:
	@mkdir -p _build
	@rm -f _build/coverage-acceptance.out
	#@cd features && LOGXI="*=ERR,dat:sqlx=OFF,dat=OFF" go test -cover -covermode=count -coverprofile=../_build/coverage-acceptance.out
	@cd features && go test -cover -covermode=count -coverprofile=../_build/coverage-acceptance.out

acceptance-run-focus:
	@mkdir -p _build
	@rm -f _build/coverage-acceptance.out
	@cd features && go test -cover -covermode=count -coverprofile=../_build/coverage-acceptance.out --focus

merge-profiles:
	@mkdir -p _build
	@gocovmerge _build/*.out > _build/coverage-all.out

test-coverage-func coverage-func: merge-profiles
	@echo
	@echo "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
	@echo "Functions NOT COVERED by Tests"
	@echo "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
	@go tool cover -func=_build/coverage-all.out | egrep -v "100.0[%]"

test-coverage-html coverage-html: merge-profiles
	@go tool cover -html=_build/coverage-all.out

run:
	@go run main.go start -v3 -c ./config/local.yaml

rtfd:
	@rm -rf docs/_build
	@sphinx-build -b html -d ./docs/_build/doctrees ./docs/source/ docs/_build/html
	@open docs/_build/html/index.html
