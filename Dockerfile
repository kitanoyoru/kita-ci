# TODO: Copy config to /etc and create sustemd service

FROM golang:latest AS build

LABEL maintainer="Alexandr Rutkovskij <kitanoyoru@protonmail.com>"

ENV PROJECT_DIR /usr/src/kita-ci

WORKDIR PROJECT_DIR

COPY . $PROJECT_DIR

RUN go mod tidy
RUN go build -o /bin/kita-ci cmd/main.go

RUN rm -rf $PROJECT_DIR


FROM debian:buster-slim

ENV USER kitanoyoru

RUN useradd -m -U -d /home/$USER $USER -s /bin/bash

RUN set -ex; \
  DEPS='vim'; \
  apt-get update; \
  apt-get install -y $DEPS --no-install-recommends;

COPY --from=build /bin/kita-ci /bin/kita-ci 

EXPOSE 8080

ENTRYPOINT ["/bin/kita-ci"]

CMD ["start"]

