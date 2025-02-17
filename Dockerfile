FROM golang:1.22.3-alpine AS build

ARG COMMIT

ENV GOPROXY "https://goproxy.io,direct"

RUN apk update && \
    apk add --no-cache tzdata ca-certificates make gettext bash

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod tidy

COPY . .

RUN make build COMMIT=${COMMIT}

FROM alpine:3.18

RUN apk update && \
    apk add --no-cache tzdata ca-certificates gettext openssl bash

RUN addgroup -S app && adduser -S -g app app

WORKDIR /app

COPY --from=build /app/config.yaml .
COPY --from=build /app/server .

RUN chmod +x /app/server

RUN chown -R app /app

USER app

EXPOSE 8080

ENTRYPOINT [ "./server" ]
