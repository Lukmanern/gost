package service

import (
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
	"sync"

	"github.com/Lukmanern/gost/internal/env"
	"github.com/Lukmanern/gost/internal/response"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type EmailService interface {
	Handler(c *fiber.Ctx) (err error)
	Send(emails []string, subject string, message string) (res map[string]bool, err error)
}

type EmailServiceImpl struct {
	Server   string
	Port     int
	Email    string
	Password string
}

var (
	emailService     *EmailServiceImpl
	emailServiceOnce sync.Once
)

func NewEmailService() EmailService {
	emailServiceOnce.Do(func() {
		config := env.Configuration()
		emailService = &EmailServiceImpl{
			Server:   config.SMTPServer,
			Port:     config.SMTPPort,
			Email:    config.SMTPEmail,
			Password: config.SMTPPassword,
		}
	})

	return emailService
}

const simpleMessage = `Lorem ipsum, dolor sit amet consectetur adipisicing elit. Ad consequuntur 
similique voluptatibus ab enim harum dolor, sit, corporis repellendus culpa cum, quasi corrupti! 
Impedit inventore cum optio quas, nisi aliquid ullam omnis voluptas, architecto deserunt, sint 
tempora? Iure ea alias recusandae sunt ad, vero laudantium esse.`

var testEmails = []string{"lukmanernandi16@gmail.com", "unsurlukman@gmail.com", "code_name_safe_in_unsafe@proton.me",
	"lukmanernandi16@gmail.com.", "unsurlukm an@gmail.com", "code _name_safe_in_unsafe@proton.me", "lukmanern*a)ndi16@gmail.com",
	"unsurlukman@gmail.com", "code_n}ame_safe_in_unsafe@proton.me",
}

func (svc EmailServiceImpl) Handler(c *fiber.Ctx) (err error) {
	res, err := svc.Send(testEmails, "Testing Gost Project", simpleMessage)
	if err != nil {
		return response.ErrorWithData(c, "internal server error: "+err.Error(), fiber.Map{
			"res": res,
		})
	}
	if res == nil {
		return response.Error(c, "internal server error: failed sending email")
	}

	message := "success sending emails"
	return response.CreateResponse(c, fiber.StatusAccepted, true, message, nil)
}

func (svc EmailServiceImpl) getAuth() smtp.Auth {
	return smtp.PlainAuth("", svc.Email, svc.Password, svc.Server)
}

func (svc EmailServiceImpl) getSMTPAddr() string {
	return fmt.Sprintf("%s:%d", svc.Server, svc.Port)
}

func (svc EmailServiceImpl) getMime() string {
	return "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
}

func (svc EmailServiceImpl) Send(emails []string, subject string, message string) (map[string]bool, error) {
	if validateErr := validateEmails(emails...); validateErr != nil {
		return nil, validateErr
	}

	subject = cases.Title(language.Und).String(subject)
	lenEmails := len(emails)
	errorSends := make([]error, lenEmails)
	var wg sync.WaitGroup

	for i, email := range emails {
		addr := svc.getSMTPAddr()
		auth := svc.getAuth()
		mime := svc.getMime()
		body := "From: " + svc.Email + "\n" +
			"To: " + email + "\n" +
			"Subject: " + subject + "\n" + mime +
			message

		wg.Add(1)
		go func(i int, email string) {
			defer func() {
				wg.Done()
			}()

			errSend := smtp.SendMail(addr, auth, svc.Email, []string{email}, []byte(body))
			if errSend != nil {
				errorSends[i] = errSend
			}
		}(i, email)
	}
	wg.Wait()

	var hasError error = nil
	res := make(map[string]bool, lenEmails)
	for i, email := range emails {
		if errorSends[i] != nil {
			res[email] = false
			hasError = errors.New("emails may have failed, check $res for detail, in $res true for success")
			continue
		}
		res[email] = true
	}

	if hasError != nil {
		return res, hasError
	}
	return res, nil
}

func validateEmails(emails ...string) error {
	for _, email := range emails {
		_, err := mail.ParseAddress(email)
		if err != nil {
			return errors.New("one or more email/s is invalid " + email)
		}
	}
	return nil
}
