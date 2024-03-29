package helper

import (
	"errors"
	"math/rand"
	"net"
	"net/mail"
	"strings"
	"time"

	"github.com/Lukmanern/gost/internal/middleware"
	"github.com/XANi/loremipsum"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func RandomWords(n int) string {
	if n < 2 {
		n = 2
	}
	loremIpsumGenerator := loremipsum.New()
	words := loremIpsumGenerator.Words(n)
	return words
}

// RandomString func generate random string
// used for testing and any needs.
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

// RandomEmail func return a email
// used for testing and any needs.
func RandomEmail() string {
	body := strings.ToLower(RandomString(7) + RandomString(7) + RandomString(7))
	randEmail := body + "@gost.project"
	return randEmail
}

// RandomIPAddress func return a IP Address
// used for testing and any needs.
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

// ValidateEmails func validates emails
func ValidateEmails(emails ...string) error {
	for _, email := range emails {
		_, err := mail.ParseAddress(email)
		if err != nil {
			return errors.New("one or more email/s is invalid: " + email)
		}
	}
	return nil
}

// NewFiberCtx func create new fiber.Ctx used for testing
// handler like controller and middleware.
func NewFiberCtx() *fiber.Ctx {
	app := fiber.New()
	return app.AcquireCtx(&fasthttp.RequestCtx{})
}

// ToTitle func make string to Title Case
// Example : Your name => Your Name
func ToTitle(s string) string {
	return cases.Title(language.Und).String(s)
}

// Generate token for admin role : Full Access
func GenerateToken() string {
	jwtHandler := middleware.NewJWTHandler()
	expire := time.Now().Add(15 * time.Hour)
	token, err := jwtHandler.GenerateJWT(GenerateRandomID(), RandomEmail(), map[string]uint8{"admin": 1}, expire)
	if err != nil {
		return ""
	}
	return token
}

func GenerateRandomID() int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	min := 9000000
	max := 10000000
	return rand.Intn(max-min) + min
}
