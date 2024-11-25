FROM golang:1.23.2 as init-stage

FROM ubuntu:latest

WORKDIR /app

COPY . .

FROM init-stage AS build-stage

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go mod download
RUN go build -o /server

ARG PORT=7540
ENV PORT=${PORT}
ARG PASSWORD=1234
ENV PASSWORD=${PASSWORD}
ARG TODODBFILE=scheduler.db
ENV TODO_DBFILE=${TODODBFILE}

EXPOSE ${PORT}
CMD ["./app"]
