FROM golang:1.21.3-alpine3.18 as builder

# gcc for CGO_ENABLED=1
RUN apk add build-base

COPY backend /build

WORKDIR /build

RUN go env -w GOPROXY=https://goproxy.cn,https://gocenter.io,https://goproxy.io,direct

RUN go build -o easyddns main.go


FROM alpine:3.18.4 as prod

WORKDIR /app

COPY --from=builder /build/easyddns .

CMD ["./easyddns"]
