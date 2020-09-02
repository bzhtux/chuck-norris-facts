# builder
FROM golang:alpine AS build-env
LABEL maintainer="Yannick Foeillet <bzhtux@gmail.com>"

RUN apk --no-cache add build-base git mercurial gcc
RUN mkdir -p /go/src/github.com/bzhtux/chuck-norris-facts
ADD . /go/src/github.com/bzhtux/chuck-norris-facts
RUN cd /go/src/github.com/bzhtux/chuck-norris-facts && go get . && go build -o cnf


# final image
FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/bzhtux/chuck-norris-facts/cnf /app/
COPY templates /app/templates
ENTRYPOINT ./cnf
