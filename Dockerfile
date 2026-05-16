FROM golang:1.25.3-alpine AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /app/main

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/static ./static
COPY --from=build /app/favi ./favi
COPY --from=build /app/templates ./templates

COPY --from=build /app/main .
ENV DOMAIN=kemcbride.noho.st
ENV GIN_MODE=release

EXPOSE 8080

CMD ["./main"]

