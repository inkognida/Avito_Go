FROM golang:latest as builder

WORKDIR /app

COPY ./ ./
RUN go build -o main ./main.go


COPY --from=builder /app/main ./

ENTRYPOINT ["/main"]
