FROM golang:alpine as build

WORKDIR /etc/source

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app

FROM scratch

WORKDIR /

COPY --from=build /app /app

EXPOSE 8090

CMD [ "/app" ]