.Phony: run

run:
	go run ./service/cmd/main.go

publish:
	go run ./stan-publisher/main.go