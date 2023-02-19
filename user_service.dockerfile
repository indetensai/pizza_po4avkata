FROM golang:alpine
WORKDIR /usr/src/app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY user_service/*.go ./
COPY user_service/.env ./
COPY service/ ./service/
RUN go build -o ./user
EXPOSE 443
CMD [ "./user" ]