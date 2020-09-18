FROM buildpack-deps:buster

ENV PERL_VERSION 5.32.0

WORKDIR /usr/local/bin

RUN curl -fsSL --compressed https://raw.githubusercontent.com/tokuhirom/Perl-Build/master/perl-build > perl-build \
    && chmod +x perl-build

RUN curl -fsSL --compressed https://raw.githubusercontent.com/skaji/cpm/master/cpm > cpm \
    && chmod +x cpm

# XXX: 緯度経度が16桁程度あっても扱いやすいように拡張倍精度浮動小数点数を利用する
RUN perl-build $PERL_VERSION /opt/perl-$PERL_VERSION/ -Duselongdouble

ENV PATH=/opt/perl-$PERL_VERSION/bin:$PATH
ENV PERL5LIB=/opt/perl-$PERL_VERSION/lib

EXPOSE 1323

WORKDIR /app

RUN apt-get update && apt-get install -y wget default-mysql-client

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

COPY cpanfile .
RUN cpm install -g --show-build-log-on-failure
COPY . .

#CMD ["plackup", "app.psgi"]
