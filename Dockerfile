# syntax=docker/dockerfile:1

FROM golang:1.19 AS build
WORKDIR /src
COPY . /src/
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/syncerd cmd/app/main.go

FROM scratch AS bin
COPY --from=build /out/syncerd /
CMD ["/syncerd", "-src", "/var/src", "-dst", "/var/dst", "-level", "info", "-format", "json", "-interval", "10s", "--force"]
