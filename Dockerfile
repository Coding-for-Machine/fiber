FROM golang:1.23 AS builder

WORKDIR /app

# Go mod va loyiha fayllarini yuklash
COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Go build qilish
RUN go build -o ./app/main 
# Yangi container yaratish
FROM debian:bullseye-slim

WORKDIR /app

# Build qilingan binarni ko‘chirish
COPY --from=builder /app/main .

# Faylga ruxsat berish
RUN chmod +x main 
 # ✅ "main" faylni bajariladigan qilib qo‘yamiz

EXPOSE 3000

CMD ["./main"] 
 # ✅ Docker "main" faylini ishga tushiradi
