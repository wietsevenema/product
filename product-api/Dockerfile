FROM golang:1.15-buster as base

WORKDIR /src
COPY go.* ./
RUN go mod download

FROM base as import
COPY assets assets/
COPY cmd/import cmd/import
COPY internal internal
RUN go run products/cmd/import

FROM base as build
COPY cmd cmd
COPY internal internal
RUN go build products/cmd/server

# Debian10 is Debian Buster (same as the build image)
FROM gcr.io/distroless/base-debian10:nonroot
WORKDIR /app
COPY --from=import /src/products.db /app/products.db
COPY --from=build /src/server /app/server

ENTRYPOINT ["./server"]