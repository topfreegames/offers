setup:
	@glide install

acceptance acc:
	@go test -tags=acceptance
