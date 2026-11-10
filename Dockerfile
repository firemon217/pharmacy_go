WORKDIR /app
# Копируем go.mod и go.sum отдельно – так кэш слоёв Docker
# переиспользуется, если изменился только код

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .
# === Этап 2: рантайм ===

FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
ENTRYPOINT ["./server"]