FROM golang:alpine3.15

WORKDIR /app/source-code

COPY source-code /app/source-code/

RUN go mod download

RUN go build -o go-crud-api

EXPOSE 8080

CMD [ "./go-crud-api" ]