package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oschwald/geoip2-golang"
)

// IP manzilini olish funksiyasi (proxy orqali ham tekshiradi)
func getIP(c *fiber.Ctx) string {
	ip := c.Get("X-Forwarded-For") // Proxy orqali kelgan IP
	if ip == "" {
		ip = c.Get("X-Real-IP") // Haqiqiy foydalanuvchi IP-si
	}
	if ip == "" {
		ip = c.IP() // Agar yuqoridagilar bo‘lmasa, oddiy IP olamiz
	}
	return ip
}

// Middleware: Foydalanuvchining IP manzili va GeoIP ma'lumotlarini logga yozish
func GeoIPLoggerMiddleware(cityDB, asnDB *geoip2.Reader) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIP := getIP(c)
		ip := net.ParseIP(userIP)

		// GeoIP2 City bo‘yicha ma'lumot olish
		cityRecord, err := cityDB.City(ip)
		if err != nil {
			fmt.Println("IP bo‘yicha shahar ma'lumotlarini olishda xatolik:", err)
		}

		// GeoIP2 ASN (ISP va Provider) bo‘yicha ma'lumot olish
		asnRecord, err := asnDB.ASN(ip)
		if err != nil {
			fmt.Println("IP bo‘yicha ASN (ISP) ma'lumotlarini olishda xatolik:", err)
		}

		// JSON log tuzish
		logEntry := map[string]interface{}{
			"time":      time.Now().Format(time.RFC3339),
			"ip":        userIP,
			"method":    c.Method(),
			"path":      c.Path(),
			"country":   cityRecord.Country.Names["en"],
			"city":      cityRecord.City.Names["en"],
			"latitude":  cityRecord.Location.Latitude,
			"longitude": cityRecord.Location.Longitude,
			"isp":       asnRecord.AutonomousSystemOrganization, // ISP nomi
			"asn":       asnRecord.AutonomousSystemNumber,       // ASN raqami
		}

		// JSON formatga o'tkazish
		logData, err := json.Marshal(logEntry)
		if err != nil {
			fmt.Println("JSON formatga o'tkazishda xatolik:", err)
			return c.Next()
		}

		// Faylga yozish
		file, err := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Fayl ochishda xatolik:", err)
			return c.Next()
		}
		defer file.Close()

		file.WriteString(string(logData) + "\n")

		// Davom etish
		return c.Next()
	}
}

func main() {
	// GeoIP bazalarini yuklash
	cityDB, err := geoip2.Open("GeoLite2-City.mmdb") // Shahar va mamlakat uchun
	if err != nil {
		log.Fatal("GeoIP City bazasi yuklanmadi:", err)
	}
	defer cityDB.Close()

	asnDB, err := geoip2.Open("GeoLite2-ASN.mmdb") // ISP va ASN uchun
	if err != nil {
		log.Fatal("GeoIP ASN bazasi yuklanmadi:", err)
	}
	defer asnDB.Close()

	// Fiber ilovasini yaratish
	app := fiber.New()

	// Middleware ulash
	app.Use(GeoIPLoggerMiddleware(cityDB, asnDB))

	// API guruhlari
	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Get("/list", handler)
	v1.Get("/user", handler)

	v2 := api.Group("/v2")
	v2.Get("/list", handler)
	v2.Get("/user", handler)

	// Serverni ishga tushirish
	log.Fatal(app.Listen(":3000"))
}

// Oddiy handler
func handler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello, World!"})
}
