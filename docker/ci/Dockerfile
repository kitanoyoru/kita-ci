FROM ubuntu:20.04

LABEL maintainer="Alexandr Rutkovskij <kitanoyoru@protonmail.com>"

ENV USER kitanoyoru

RUN useradd -m -U -d /home/$USER $USER -s /bin/bash

RUN set -ex; \
  DEPS='git curl'; \
  apt-get update; \
  apt-get install -y $DEPS --no-install-recommends;

ENV DOCKER_VERSION=20.10.9

RUN curl -fsSLO https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKERVERSION}.tgz \
  && tar xzvf docker-${DOCKERVERSION}.tgz --strip 1 \
                 -C /usr/local/bin docker/docker \
  && rm docker-${DOCKERVERSION}.tgz


COPY docker/docker-entrypoint.sh /usr/local/bin/
RUN chmod 755 /usr/local/bin/docker-entrypoint.sh

USER $USER

RUN mkdir /builder
WORKDIR /builder

ENTRYPOINT ["docker-entrypoint.sh"]
