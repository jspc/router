FROM alpine
MAINTAINER jspc

RUN apk add --update --no-cache ca-certificates

ADD router /bin

ENTRYPOINT ["/bin/router"]
