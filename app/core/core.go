package core

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type AppCore struct {
	DB *sqlx.DB
	tz *time.Location
}

func New(db *sqlx.DB) *AppCore {
	return &AppCore{
		DB: db,
	}
}

// func New(db *sqlx.DB, appApp app.App) *Core {
// 	appInfo := appApp.AppInfo()

// 	inst := Core{
// 		db:  db,
// 		App: appApp,
// 		tz:  appInfo.TZ(),
// 	}
// 	clientBase, err := gatewayrest.NewServiceClient(db)
// 	if err != nil {
// 		panic(err)
// 	}
// 	inst.ServiceClient = clientBase
// 	return &inst
// }

// func (core *Core) DB() *gorm.DB {
// 	return core.db
// }
