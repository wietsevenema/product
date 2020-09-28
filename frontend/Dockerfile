# A named stage: "build"
FROM golang:1.15 AS builder
 
WORKDIR /go/src/app
COPY go.* /go/src/app/
RUN go mod download
 
COPY . /go/src/app
RUN go build -o /go/bin/app

# A new stage: "run"
FROM gcr.io/distroless/base-debian10:nonroot AS run

# Copy the binary from stage build
COPY --from=builder /go/bin/app .
COPY --chown=nonroot:nonroot web web
 
ENTRYPOINT ["./app"]
