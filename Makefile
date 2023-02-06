RUN = @go run cmd/app/main.go

all: run

dry:
	$(RUN)

run:
	$(RUN) -src ./tmp/src -dst ./tmp/dst -log log.json -level info -format json -interval 10s $(ARGS)

run_force:
	@make run ARGS="--force"

run_jq:
	@make run | jq .

run_warn:
	@make run ARGS="-level warning"

test:
	@go test ./... -cover -count=1 -coverprofile=cover.out

bench:
	go test -test.bench BenchmarkCopy ./pkg/syncd/ -benchtime=100x

show_cover:
	@go tool cover -html=cover.out

cover: test show_cover

build:
	@go build -o syncerd cmd/app/main.go
