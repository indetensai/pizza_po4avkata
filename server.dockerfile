FROM golang:latest
WORKDIR /usr/src/app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY server/*.go ./
COPY service ./service/
RUN go build -o ./server
EXPOSE 8080
CMD [ "./server" ]