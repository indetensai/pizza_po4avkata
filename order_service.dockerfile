FROM golang:latest
WORKDIR /usr/src/app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY order_service/*.go ./
COPY order_service/.env ./
COPY service/ ./service/
RUN go build -o ./order
EXPOSE 50000
CMD [ "./order" ]