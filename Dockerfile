# Go bazaviy imijidan foydalanamiz
FROM golang:1.21-alpine

# Ishlash katalogini o‘rnatamiz
WORKDIR /app

# Kerakli paketlarni o‘rnatamiz
RUN apk add --no-cache curl unzip

# GeoLite2 City bazasini yuklash
RUN curl -L -o GeoLite2-City.mmdb.gz "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb.gz" \
    && gunzip GeoLite2-City.mmdb.gz

# GeoLite2 ASN bazasini yuklash
RUN curl -L -o GeoLite2-ASN.mmdb.gz "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-ASN.mmdb.gz" \
    && gunzip GeoLite2-ASN.mmdb.gz

# Loyihani ko‘chiramiz
COPY . .

# Modullarni yuklash
RUN go mod tidy

# Loyihani build qilish
RUN go build -o app /app/main.go

# Port ochamiz
EXPOSE 3000

# Ilovani ishga tushiramiz
CMD [ "./app" ]
