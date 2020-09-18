FROM rust:1.46

EXPOSE 1323

WORKDIR /usr/src/isuumo

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

COPY Cargo.toml Cargo.lock ./
RUN mkdir src && echo 'fn main() {}' > src/main.rs && cargo build --locked --release
COPY src ./src
RUN touch src/main.rs && cargo build --locked --frozen --release

RUN apt-get update && apt-get install -y wget default-mysql-client

CMD ["/usr/src/isuumo/target/release/isuumo"]
