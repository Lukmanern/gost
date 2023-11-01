package controller

import (
	"sync"
	"time"

	"github.com/Lukmanern/gost/database/connector"
	"github.com/Lukmanern/gost/internal/response"
	"github.com/gofiber/fiber/v2"
)

type DevController interface {
	// Database
	PingDatabase(c *fiber.Ctx) error
	PingRedis(c *fiber.Ctx) error
	// Handler
	Panic(c *fiber.Ctx) error
	// Redis
	StoringToRedis(c *fiber.Ctx) error
	GetFromRedis(c *fiber.Ctx) error
}

type DevControllerImpl struct{}

var (
	devImpl     *DevControllerImpl
	devImplOnce sync.Once
)

func NewDevControllerImpl() DevController {
	devImplOnce.Do(func() {
		devImpl = &DevControllerImpl{}
	})

	return devImpl
}

func (ctr DevControllerImpl) PingDatabase(c *fiber.Ctx) error {
	db := connector.LoadDatabase()
	sqldb, sqlErr := db.DB()
	if sqlErr != nil {
		return response.Error(c, "failed get sql-db")
	}
	for i := 0; i < 5; i++ {
		pingErr := sqldb.Ping()
		if pingErr != nil {
			return response.Error(c, "failed to ping-sql-db")
		}
	}

	return response.CreateResponse(c, fiber.StatusOK, true, "success ping-sql-db", nil)
}

func (ctr DevControllerImpl) PingRedis(c *fiber.Ctx) error {
	redis := connector.LoadRedisDatabase()
	for i := 0; i < 5; i++ {
		status := redis.Ping()
		if status.Err() != nil {
			return response.Error(c, "failed to ping-redis")
		}
	}

	return response.CreateResponse(c, fiber.StatusOK, true, "success ping-redis", nil)
}

func (ctr DevControllerImpl) Panic(c *fiber.Ctx) error {
	defer func() error {
		r := recover()
		if r != nil {
			return response.Error(c, "message panic: "+r.(string))
		}
		return nil
	}()
	panic("Panic message") // message should string
}

func (ctr DevControllerImpl) StoringToRedis(c *fiber.Ctx) error {
	redis := connector.LoadRedisDatabase()
	if redis == nil {
		return response.Error(c, "redis nil value")
	}
	redisStatus := redis.Set("example-key", "example-value", 50*time.Minute)
	if redisStatus.Err() != nil {
		return response.Error(c, "redis status error ("+redisStatus.Err().Error()+")")
	}

	return response.SuccessCreated(c, nil)
}

func (ctr DevControllerImpl) GetFromRedis(c *fiber.Ctx) error {
	redis := connector.LoadRedisDatabase()
	if redis == nil {
		return response.Error(c, "redis nil value")
	}
	redisStatus := redis.Get("example-key")
	if redisStatus.Err() != nil {
		return response.Error(c, "redis status error ("+redisStatus.Err().Error()+")")
	}
	res, resErr := redisStatus.Result()
	if resErr != nil {
		return response.Error(c, "redis result error ("+resErr.Error()+")")
	}

	return response.SuccessLoaded(c, res)
}
