// 📌 Origin Github Repository: https://github.com/Lukmanern

package application

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/Lukmanern/gost/database/connector"
	"github.com/Lukmanern/gost/internal/env"
)

var (
	port int

	// Create a new fiber instance with custom config
	router = fiber.New(fiber.Config{
		AppName: "Gost Project",
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code
			// if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Send custom error page
			err = ctx.Status(code).JSON(fiber.Map{
				"message": e.Message,
			})
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).
					SendString("Internal Server Error")
			}
			return nil
		},
		// memory management
		// ReduceMemoryUsage: true,
		// ReadBufferSize: 5120,
	})
)

// setup initializes the application
// by checking the environment and
// database configuration.
func setup() {
	env.ReadConfig("./.env")
	config := env.Configuration()
	privKey := config.GetPrivateKey()
	pubKey := config.GetPublicKey()
	if privKey == nil || pubKey == nil {
		log.Fatal("private and public keys are not valid or not found")
	}
	port = config.AppPort

	connector.LoadDatabase()
	connector.LoadRedisCache()
}

// checkLocalPort checks if a given port is available for local use.
func checkLocalPort(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("port is being used by other process")
	}
	defer listener.Close()
}

// RunApp initializes and runs the application,
// handling setup, port checking, middleware,
// and route registration.
func RunApp() {
	setup()
	checkLocalPort(port)
	router.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	router.Use(logger.New())
	// Custom File Writer
	_ = os.MkdirAll("./log", os.ModePerm)
	fileName := fmt.Sprintf("./log/%s.log", time.Now().Format("20060102"))
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	router.Use(logger.New(logger.Config{
		Output: file,
	}))

	// Create channel for idle connections.
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		// Received an interrupt signal, shutdown.
		// ctrl+c
		if err := router.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnsClosed)
	}()

	helloRoutes(router)
	userRoutes(router)
	roleRoutes(router)

	if err := router.Listen(fmt.Sprintf(":%d", port)); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}
