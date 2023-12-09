build_and_run: build run

build:
	go build
run:
	./smalltown
test:
	go test -coverprofile cover.out
show_coverage:
	go tool cover -func=cover.out
	go tool cover -html=cover.out
