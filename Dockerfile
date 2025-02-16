# Go bazaviy imijidan foydalanamiz
FROM golang:1.23 AS builder 

# Ishlash katalogini o‘rnatamiz
WORKDIR /app


# Loyihani ko‘chiramiz
COPY . .

# Modullarni yuklash
RUN go mod tidy

# Loyihani build qilish
RUN go build -o app ./app/main.go

# Port ochamiz
EXPOSE 3000

# Ilovani ishga tushiramiz
CMD [ "./app" ]
