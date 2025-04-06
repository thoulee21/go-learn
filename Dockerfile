FROM golang:1.24-alpine AS builder
WORKDIR /srv/go-app
COPY . .

RUN apk add --no-cache gcc g++ make autoconf2.13 automake1.15 binutils-dev gmp-dev isl-dev libmpc-dev libgcc-dev patchelf

RUN go env -w CGO_ENABLED=1
RUN go build -o aichatbot


FROM alpine
WORKDIR /srv/go-app
COPY --from=builder /srv/go-app/aichatbot .

CMD ["./aichatbot"]