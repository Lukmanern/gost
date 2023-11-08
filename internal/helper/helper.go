package helper

import (
	"errors"
	"math/rand"
	"net"
	"net/mail"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func RandomString(n uint) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyz"
	letterBytes += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterBytes += "1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RandomEmails(n uint) []string {
	emailsMap := make(map[string]int)
	for uint(len(emailsMap)) < n {
		body := strings.ToLower(RandomString(7) + RandomString(7) + RandomString(7))
		randEmail := body + "@gost.project"
		emailsMap[randEmail] += 1
	}

	emails := make([]string, 0, len(emailsMap))
	for email := range emailsMap {
		emails = append(emails, email)
	}
	return emails
}

func RandomEmail() string {
	body := strings.ToLower(RandomString(7) + RandomString(7) + RandomString(7))
	randEmail := body + "@gost.project"
	return randEmail
}

func RandomIPAddress() string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	a := rng.Intn(256)
	b := rng.Intn(256)
	c := rng.Intn(256)
	d := rng.Intn(256)
	ip := net.IPv4(byte(a), byte(b), byte(c), byte(d))

	return ip.String()
}

func ValidateEmails(emails ...string) error {
	for _, email := range emails {
		_, err := mail.ParseAddress(email)
		if err != nil {
			return errors.New("one or more email/s is invalid " + email)
		}
	}
	return nil
}

// This used for testing handler : controller/ middleware/ any
func NewFiberCtx() *fiber.Ctx {
	app := fiber.New()
	return app.AcquireCtx(&fasthttp.RequestCtx{})
}

func ToTitle(s string) string {
	return cases.Title(language.Und).String(s)
}
