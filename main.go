package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hublabs/common/auth"
	"github.com/hublabs/product-availability-api/config"
	"github.com/hublabs/product-availability-api/controllers"
	"github.com/hublabs/product-availability-api/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pangpanglabs/echoswagger"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/sirupsen/logrus"
)

func main() {
	config := config.Init(os.Getenv("APP_ENV"))

	db := initDB(config.Database.Stock.Driver, config.Database.Stock.Connection)
	models.Init(db)
	migrate(db)

	e := echo.New()

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	e.GET("/swagger", func(c echo.Context) error {
		return c.File("./swagger.yml")
	})
	e.GET("/whoami", func(c echo.Context) error {
		return c.String(http.StatusOK, config.ServiceName)
	})
	e.Static("/docs", "./swagger-ui")

	r := echoswagger.New(e, "/docs", &echoswagger.Info{
		Title:       "Echo Sample",
		Description: "This is API doc for Echo Sample",
		Version:     "1.0",
	}).AddSecurityAPIKey("Authorization", "JWT Token", "header")

	controllers.StockController{}.Init(r.Group("Stocks", "/v1/stocks"))

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	// e.Use(middleware.RequestID())
	e.Use(echomiddleware.BehaviorLogger(config.ServiceName, config.BehaviorLog.Kafka))
	e.Use(echomiddleware.ContextDB(config.ServiceName, db, config.Database.Logger.Kafka))
	e.Use(auth.UserClaimMiddleware("/ping", "/docs"))

	if config.AppEnv != "production" {
		behaviorlog.SetLogLevel(logrus.InfoLevel)
		logrus.SetLevel(logrus.InfoLevel)
	}

	if err := e.Start(":8000"); err != nil {
		log.Println(err)
	}

}
func initDB(driver, connection string) *xorm.Engine {
	db, err := xorm.NewEngine(driver, connection)
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(time.Minute * 10)
	db.ShowSQL()

	return db
}
func migrate(db *xorm.Engine) {
	if exist, err := existColumn(db, "stock_distribution", "status"); err != nil {
		logrus.WithError(err).Warning("Fail to migrate")
	} else if exist {
		if _, err := db.Exec("ALTER TABLE `stock_distribution` DROP COLUMN `status`"); err != nil {
			logrus.WithError(err).Info("Fail to migrate")
		}
	}

	if exist, err := existColumn(db, "stock_distribution", "stock_type"); err != nil {
		logrus.WithError(err).Warning("Fail to migrate")
	} else if exist {
		if _, err := db.Exec("ALTER TABLE `stock_distribution` DROP COLUMN `stock_type`"); err != nil {
			logrus.WithError(err).Info("Fail to migrate")
		}
	}
}

func existColumn(db *xorm.Engine, tableName, columnName string) (bool, error) {
	rows, err := db.Query("SELECT COLUMN_NAME FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ? AND `COLUMN_NAME` = ?",
		"omni_stock", tableName, columnName)
	if err != nil {
		return false, nil
	}
	return len(rows) == 1, nil
}
