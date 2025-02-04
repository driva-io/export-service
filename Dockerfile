ARG GO_VERSION=1.23.0
FROM golang:${GO_VERSION} AS build

WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY . .

ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod/ \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/consumer ./cmd/consumer

RUN --mount=type=cache,target=/go/pkg/mod/ \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/http ./cmd/http

FROM alpine:latest AS final

RUN apk --no-cache add \
        ca-certificates \
        tzdata \
        && update-ca-certificates

COPY --from=build /bin/consumer /bin/
COPY --from=build /bin/http /bin/
COPY start.sh /bin/

EXPOSE 23545

# Run start.sh
CMD sh /bin/start.sh
