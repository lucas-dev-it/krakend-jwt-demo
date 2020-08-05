FROM golang:1.13.4  as builder
WORKDIR /go/src/bitbucket.org/rappinc/cpgs-int-worker
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/bitbucket.org/rappinc/cpgs-int-worker/app .
CMD ["./app"]