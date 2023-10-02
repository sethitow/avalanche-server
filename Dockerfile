FROM golang:1.21-bookworm

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY ./ ./

RUN GOOS=linux go build ./cmd/avalancheserver

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=0 /app/avalancheserver /bin/avalancheserver

EXPOSE 8080

CMD ["avalancheserver"]
