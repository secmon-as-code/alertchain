FROM golang:1.20.2 AS build-go
ENV CGO_ENABLED=0
ARG BUILD_VERSION
COPY . /app
WORKDIR /app
RUN go build -o alertchain -ldflags "-X github.com/m-mizutani/alertchain/pkg/domain/types.AppVersion=${BUILD_VERSION}" .

FROM gcr.io/distroless/base
COPY --from=build-go /app/alertchain /alertchain

ENTRYPOINT ["/alertchain"]
