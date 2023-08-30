FROM golang:1.20-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o /app/main .

FROM scratch
COPY --from=builder /app/ /app

EXPOSE 8090
ENTRYPOINT [ "/app/main" ]