FROM alpine:3.7

LABEL maintainer=i@tosiney.com

ARG GOPATH=/data/gocode
ARG CONFIG=config.yml

WORKDIR /data

COPY ${CONFIG} /data

RUN apk add --no-cache --virtual .build-deps \
  gcc go make libc-dev wget curl \
  && wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh \
  && chmod +x wait-for-it.sh \
  && VERSION=`curl --silent "https://api.github.com/repos/EffDataAly/GithubTraveler/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'` \
  && mkdir -p /data/gocode/src/github.com/EffDataAly \
  && cd /data/gocode/src/github.com/EffDataAly \
  && wget https://github.com/EffDataAly/GithubTraveler/archive/${VERSION}.tar.gz \
  && tar -zxvf ${VERSION}.tar.gz \
  && mv GithubTraveler-${VERSION} GithubTraveler \
  && cd GithubTraveler \
  && make \
  && apk del .build-deps \
  && cp ./release/GithubTraveler-linux /usr/bin \
  && cd /data \
  && rm -rf /data/gocode \
  && apk add --no-cache ca-certificates bash

VOLUME /data/db

CMD ["GithubTraveler-linux", "crawler", "--config", "config.yml"]
