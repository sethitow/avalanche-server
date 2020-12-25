FROM rust:1.48 as builder

WORKDIR ./avalanche-server
RUN rustup default nightly

COPY . ./

RUN cargo build --release


FROM debian:buster-slim

COPY --from=builder /avalanche-server/target/release/avalanche_server .
COPY avalanche_data.json .

EXPOSE 8000
CMD ["./avalanche_server"]