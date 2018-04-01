FROM alpine

WORKDIR /data

ARG GOPATH=/data/gocode
ARG CONFIG=config.yml

COPY ${CONFIG} /data

RUN apk add --no-cache --virtual .build-deps \
  gcc go git make libc-dev wget \
  && wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh \
  && mkdir -p /data/gocode/src/github.com/EffDataAly \
  && cd /data/gocode/src/github.com/EffDataAly \
  && git clone https://github.com/EffDataAly/GithubTraveler.git \
  && cd GithubTraveler \
  && make linux \
  && apk del .build-deps \
  && cp ./release/GithubTraveler-linux /usr/bin \
  && cd /data \
  && rm -rf /data/gocode \
  && apk add --no-cache ca-certificates

VOLUME /data/db

CMD ["GithubTraveler-linux", "crawler", "--config", "config.yml"]
