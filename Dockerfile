FROM golang:1.24-alpine AS builder
WORKDIR /srv/go-app
COPY . .

RUN apk add --no-cache gcc g++ make autoconf automake binutils-dev gmp-dev isl-dev patchelf mpc1-dev

RUN go env -w CGO_ENABLED=1
RUN go build -o aichatbot


FROM alpine
WORKDIR /srv/go-app
COPY --from=builder /srv/go-app/aichatbot .

CMD ["./aichatbot"]