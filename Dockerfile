FROM golang:1.25.5-alpine AS builder

# ARG HTTP_PROXY
# ARG HTTPS_PROXY
# ARG http_proxy
# ARG https_proxy
# ARG NO_PROXY
#
# ENV HTTP_PROXY=$HTTP_PROXY
# ENV HTTPS_PROXY=$HTTPS_PROXY
# ENV http_proxy=$http_proxy
# ENV https_proxy=$https_proxy
# ENV NO_PROXY=$NO_PROXY

RUN apk add --no-cache

WORKDIR /app

# STEP 2: Set GOPROXY WITHOUT your local proxy
RUN go env -w GOPROXY="https://proxy.golang.org,direct"

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s " -o main .
FROM alpine:3.21

RUN apk add --no-cache 
WORKDIR /app

COPY --from=builder /app/main .

COPY .env .

EXPOSE 8081

CMD ["./main"]
