
#build stage
FROM golang:alpine AS builder
COPY . /go/src/github.com/stijnv1/golang-azure/
RUN apk add --no-cache git
WORKDIR /go/src/github.com/stijnv1/golang-azure/
RUN go get -d -v ./...
RUN go install -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/ /app
ENTRYPOINT ./app/getazurevmlist
LABEL Name=golang-azure Version=0.0.1
EXPOSE 8000
