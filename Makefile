setup:
	@go get github.com/DATA-DOG/godog/cmd/godog
	@go get github.com/DATA-DOG/godog
	@glide install

acceptance acc:
	@go test -tags=acceptance
