FROM golang:latest

RUN mkdir /app
WORKDIR /app

COPY . ./    

RUN go mod download

RUN go build -o /build

RUN mkdir /data

EXPOSE 8080

ENTRYPOINT [ "/build" ]