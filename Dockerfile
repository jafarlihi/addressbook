FROM golang:1.14
WORKDIR /go/src/app
COPY . .
RUN go build

FROM debian:10.4
COPY --from=0 /go/src/app/. .
ENV WAIT_VERSION 2.7.2
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x ./wait
CMD ["./addressbook"]
