ARG GOLANG_IMAGE=""
ARG APP_BASE_IMAGE=""
FROM ${GOLANG_IMAGE} AS builder

# we copy go binary from debian based image to alpine based image so CGO must be disabled
ARG CGO_ENABLED="0"
ENV CGO_ENABLED=${CGO_ENABLED}

WORKDIR /go/src/app

RUN apt-get update && apt-get install -y \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*


RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY . .

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target="/root/.cache/go-build" \
    go build -o ./bin/api ./cmd/api/main.go

FROM ${APP_BASE_IMAGE}
# Update all distro packages to fix vulnerabilities:
RUN apk --update-cache upgrade --no-cache --available && apk add --upgrade libcrypto3 libssl3

COPY --from=builder /go/src/app/bin/api .

