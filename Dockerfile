FROM golang:1.22

WORKDIR /

COPY ./ ./

RUN go mod download
RUN go install