FROM golang:1.15-alpine3.12 as BUILD
WORKDIR /opt/print-simple
COPY . .
RUN apk add git 
RUN go get -d -v ./...
RUN go build -o print-simple

FROM alpine:3.13 as FINAL
RUN apk add inotify-tools
COPY --from=BUILD /opt/print-simple/print-simple /bin/
EXPOSE 8080
CMD ["print-simple"]