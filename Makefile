setup:
	@glide install

acceptance acc:
	@go test $(go list ./... | grep -v /vendor/) -tags=acceptance
