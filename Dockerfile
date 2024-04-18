# FROM golang:alpine

# ENV PROJECT_DIR=/app \
#     GO111MODULE=on \
#     CGO_ENABLED=0

# WORKDIR /app

# RUN mkdir "/build"

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# RUN go get github.com/githubnemo/CompileDaemon

# RUN go install github.com/githubnemo/CompileDaemon

# EXPOSE 8080

# ENTRYPOINT CompileDaemon -build="go build -o /build/app" -command="/build/app"

# -----------------------

FROM golang:alpine AS builder

ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src

COPY go.mod .
COPY go.sum .
RUN go mod download

RUN apk -U add ca-certificates
RUN apk update && apk upgrade && apk add pkgconf git bash build-base sudo
RUN git clone https://github.com/edenhill/librdkafka.git && cd librdkafka && ./configure --prefix /usr && make && make install

COPY . .

RUN go build -tags musl --ldflags "-extldflags -static" -o main .

FROM scratch AS runner

COPY --from=builder /go/src/main /

EXPOSE 8080

ENTRYPOINT ["./main"]