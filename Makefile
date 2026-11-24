server:
	@trap 'true' INT; go run ./cmd/server/; exit 0

gui:
	go run ./cmd/gui

cli:
	@trap 'true' INT; go run ./cmd/cli/; exit 0

kill:
	-fuser -k 8080/tcp


# clean:



# lint:


# lint-strict:


.PHONY: server gui cli kill
