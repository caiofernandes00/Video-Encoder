package database

import (
	"encoder/application/utils"
	"encoder/domain"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
)

type Database struct {
	Db            *gorm.DB
	Dsn           string
	DsnTest       string
	DbType        string
	DbTypeTest    string
	Debug         bool
	AutoMigrateDb bool
	Env           string
}

func NewDb() *Database {
	return &Database{}
}

func NewDbTest() *gorm.DB {
	debug, err := strconv.ParseBool(os.Getenv(utils.Debug))
	if err != nil {
		log.Fatalf("failed to parse DEBUG env to boolean")
	}

	autoMigrate, err := strconv.ParseBool(os.Getenv(utils.AutoMigrateDb))
	if err != nil {
		log.Fatalf("failed to parse AUTO_MIGRATE_DB env variable to boolean")
	}

	dbInstance := NewDb()
	dbInstance.Env = os.Getenv(utils.Env)
	dbInstance.DbTypeTest = os.Getenv(utils.DbTypeTest)
	dbInstance.DsnTest = os.Getenv(utils.DbType)
	dbInstance.AutoMigrateDb = autoMigrate
	dbInstance.Debug = debug

	connection, err := dbInstance.Connect()

	if err != nil {
		log.Fatalf("Test db error: %v", err)
	}

	return connection
}

func (d *Database) Connect() (*gorm.DB, error) {
	var err error

	if d.Env != "test" {
		d.Db, err = gorm.Open(d.DbType, d.Dsn)
	} else {
		d.Db, err = gorm.Open(d.DbTypeTest, d.DsnTest)
	}

	if err != nil {
		return nil, err
	}

	if d.Debug {
		d.Db.LogMode(true)
	}

	if d.AutoMigrateDb {
		d.Db.AutoMigrate(&domain.Video{}, &domain.Job{})
		d.Db.Model(domain.Job{}).AddForeignKey("video_id", "videos (id)", "CASCADE", "CASCADE")
	}

	return d.Db, nil
}
