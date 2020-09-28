FROM golang:1.13-buster as gobuilder

WORKDIR /app
COPY . .

ENV GOOS linux
ENV GOARCH amd64
# Create a sqllite database file with products from the products.json.gz
RUN go run products/import

# Compile the product service
RUN go build -o main products/server

# Debian10 is Debian Buster (same as the build image)
FROM gcr.io/distroless/base-debian10
WORKDIR /app
COPY --from=gobuilder /app/main /app/main
COPY --from=gobuilder /app/products.db /app/products.db

ENTRYPOINT ["./main"]