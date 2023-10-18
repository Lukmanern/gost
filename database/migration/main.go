package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Lukmanern/gost/database/connector"
	"github.com/Lukmanern/gost/domain/entity"
	"github.com/Lukmanern/gost/internal/env"
	"github.com/Lukmanern/gost/internal/rbac"
)

// Becareful using this
// This will delete entire DB Tables,
// and recreate from beginning
func main() {
	db := connector.LoadDatabase()
	fmt.Print("\n\nStart Migration\n\n")
	defer fmt.Print("\n\nFinish Migration\n\n")

	// do in development
	// Becoreful, delete entire
	// Tables and datas of Your Database.
	env.ReadConfig("./.env")
	config := env.Configuration()
	appInProduction := config.GetAppInProduction()
	if !appInProduction {
		func() {
			fmt.Print("\n\nWarning : DROPING ALL DB-TABLES AND RE-CREATE in 9 seconds (CTRL+C to stop)\n\n")
			time.Sleep(9 * time.Second)
			tables := entity.AllTables()
			deleteErr := db.Migrator().DropTable(tables...)
			if deleteErr != nil {
				log.Panicf("Error while deleting tables DB : %s", deleteErr)
			}
		}()
	}

	migrateErr := db.AutoMigrate(
		entity.AllTables()...,
	)
	if migrateErr != nil {
		log.Panicf("Error while migration DB : %s", migrateErr)
		db.Rollback()
	}

	if !appInProduction {
		// seeding permission and role
		for _, data := range rbac.AllRoles() {
			time.Sleep(100 * time.Millisecond)
			if createErr := db.Create(&data).Error; createErr != nil {
				log.Panicf("Error while create Roles : %s", createErr)
			}
		}
		time.Sleep(500 * time.Millisecond)
		for _, data := range rbac.AllPermissions() {
			time.Sleep(100 * time.Millisecond)
			if createErr := db.Create(&data).Error; createErr != nil {
				log.Panicf("Error while create Permissions : %s", createErr)
			}
			if createErr := db.Create(&entity.RoleHasPermission{
				RoleID:       1,
				PermissionID: data.ID,
			}).Error; createErr != nil {
				log.Panicf("Error while create Roles : %s", createErr)
			}
		}
	}
}
