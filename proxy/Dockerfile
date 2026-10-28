FROM golang:1.27 AS build

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -o /proxy ./cmd/proxy/main.go

FROM gcr.io/distroless/static-debian12

COPY --from=build /proxy /proxy
COPY config.yaml /config.yaml

ENV CONFIG=/config.yaml
EXPOSE 8080
ENTRYPOINT ["/proxy"]
