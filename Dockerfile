ARG GO_VERSION=1.20
FROM golang:${GO_VERSION}-alpine AS build
ARG APP_NAME=app
RUN apk update --no-cache
RUN apk add --no-cache \
  ca-certificates \
  g++ \
  git \
  make
RUN go install github.com/swaggo/swag/cmd/swag@master
WORKDIR /go/src/${APP_NAME}
COPY ./go.mod ./go.sum ./
RUN go mod download -x
COPY ./.git ./.git
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./Makefile .
ENV CGO_ENABLED=0
RUN make docs-swaggo
RUN make deps
RUN make binary
RUN sha256sum ./bin/${APP_NAME} > ./bin/${APP_NAME}.sha256

FROM scratch AS final
ARG APP_NAME=app
ENV APP_NAME=${APP_NAME}
COPY --from=build /go/src/${APP_NAME}/bin/${APP_NAME} /entrypoint
COPY --from=build /go/src/${APP_NAME}/bin/${APP_NAME}.sha256 /${APP_NAME}.sha256
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# CMD [ /bin/bash ]
ENTRYPOINT [ "/entrypoint" ]
