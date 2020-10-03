FROM golang:1.15.2

ARG WALRUS_VERSION=0.1.0

ENV GO111MODULE=on

RUN mkdir -p /app/configs
RUN mkdir -p /app/var/logs
RUN mkdir -p /app/var/storage
RUN apt-get update

WORKDIR /app

RUN curl -sL https://github.com/Clivern/Walrus/releases/download/v${WALRUS_VERSION}/Beetle_${WALRUS_VERSION}_Linux_x86_64.tar.gz | tar xz
RUN rm LICENSE
RUN rm README.md
RUN mv Walrus walrus

COPY ./config.dist.yml /app/configs/

EXPOSE 8000

VOLUME /app/configs
VOLUME /app/var

RUN ./walrus version

CMD ["./walrus", "tower", "-c", "/app/configs/config.dist.yml"]