FROM golang:latest

RUN mkdir /app
WORKDIR /app

COPY . ./    

RUN go mod download

RUN go build -o /build

EXPOSE 8080 8000

ENTRYPOINT [ "/build" ]