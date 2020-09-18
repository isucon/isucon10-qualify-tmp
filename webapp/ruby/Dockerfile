FROM rubylang/ruby:2.7.0-bionic

RUN apt-get update && \
    apt-get install -y wget build-essential default-mysql-client default-libmysqlclient-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir /app
COPY Gemfile Gemfile.lock /tmp/
RUN cd /tmp && \
  bundle config set deployment true && \
  bundle config set path /gems && \
  bundle config set without 'development test' && \
  bundle install -j4

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

WORKDIR /app
COPY . /app
EXPOSE 1323
