FROM golang:1.15-alpine AS base

WORKDIR /app

ENV PORT=$PORT
EXPOSE $PORT

LABEL maintainer="Leonardo Schettini <leoschettini2@gmail.com>"

FROM base AS builder-base

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

COPY go.mod go.sum ./

RUN go mod download

COPY . .

FROM builder-base AS development

RUN go get github.com/cespare/reflex

CMD ["reflex", "--start-service", "-r", "\\.go$", "go", "run", "./cmd/gaivota"]

FROM builder-base AS builder

RUN go build -o gaivota ./cmd/gaivota

FROM base AS production

COPY --from=builder /app/gaivota ./

CMD ["./gaivota"]
