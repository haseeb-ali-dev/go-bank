build:
	@go build -o bin/gobank

run: build
	@./bin/gobank

seed: build
	@./bin/gobank --seed

test:
	@go test -v ./...


commit:
	@if [ -z "$(msg)" ]; then \
		echo "Please provide a commit message. Usage: make commit msg='Your message'"; \
		exit 1; \
	fi
	@git add .
	@git ci -m "$(msg)"