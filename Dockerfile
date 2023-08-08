ARG GO_VERSION=1.20
FROM golang:${GO_VERSION}-alpine AS build
RUN apk add --no-cache make g++ ca-certificates
WORKDIR /go/src/app
COPY ./go.mod ./go.sum ./
RUN go mod download -x
COPY . .
ENV CGO_ENABLED=0
RUN go install github.com/swaggo/swag/cmd/swag@master
RUN make docs-swaggo
RUN make deps
RUN make binary
RUN sha256sum ./bin/app > ./bin/app.sha256

FROM scratch AS final
COPY --from=build /go/src/app/bin/app /app
COPY --from=build /go/src/app/bin/app.sha256 /app.sha256
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/app"]
