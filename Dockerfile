FROM golang:1.13 as build
ADD . /go/src/github.com/previousnext/acquia-cli
WORKDIR /go/src/github.com/previousnext/acquia-cli
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/acquia-cli /go/src/github.com/previousnext/acquia-cli/cmd/acquia-cli

FROM alpine:3.11
RUN apk --no-cache add bash ca-certificates
COPY --from=build /go/src/github.com/previousnext/acquia-cli/bin/acquia-cli /usr/local/bin/acquia-cli
WORKDIR /workspace
ENTRYPOINT ["acquia-cli"]
