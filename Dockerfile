FROM alpine

RUN apk add --no-cache --virtual .build-deps \
  gcc go git make \
  && cd /root \
  && export GOPATH=/root/gocode \
  && mkdir -p /root/gocode/src/github.com/EffDataAly \
  && cd /root/gocode/src/github.com/EffDataAly \
  && git clone https://github.com/EffDataAly/GithubTraveler.git \
  && cd GithubTraveler \
  && make linux \
  && apk del .build-deps
