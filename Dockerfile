ARG GO_VERSION=1.23.0
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build

WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY . .

ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod/ \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/app ./cmd/consumer

FROM alpine:latest AS final

RUN apk --no-cache add \
        ca-certificates \
        tzdata \
        && update-ca-certificates

COPY --from=build /bin/app /bin/

CMD [ "/bin/app" ]
