FROM alpine

WORKDIR /data

COPY config.yaml

RUN apk add --no-cache --virtual .build-deps \
  gcc go git make \
  && cd /data \
  && export GOPATH=/data/gocode \
  && mkdir -p /data/gocode/src/github.com/EffDataAly \
  && cd /data/gocode/src/github.com/EffDataAly \
  && git clone https://github.com/EffDataAly/GithubTraveler.git \
  && cd GithubTraveler \
  && make linux \
  && apk del .build-deps \
  && cp ./release/GithubTraveler-linux /usr/bin \
  && cd /data \
  && rm -rf /data/gocode

CMD ["GithubTraveler-linux", "--config", "config.yaml"]