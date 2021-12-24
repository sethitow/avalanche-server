FROM rust:1.57 as builder

WORKDIR ./avalanche-server
RUN rustup default nightly

COPY . ./

RUN cargo build --release


FROM debian:buster-slim
RUN apt-get update; apt-get install -y openssl ca-certificates
COPY --from=builder /avalanche-server/target/release/avalanche_server .

EXPOSE 8000
CMD ["./avalanche_server"]