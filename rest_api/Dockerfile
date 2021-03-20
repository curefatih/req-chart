FROM golang:1.14-alpine as builder

ENV GO111MODULE=on

WORKDIR /

COPY . .

RUN apk --no-cache add git alpine-sdk build-base gcc

RUN go build -tags musl -o ./bin/main 

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /

COPY --from=builder /bin/main .

RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]