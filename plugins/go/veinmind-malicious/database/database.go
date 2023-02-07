package database

import (
	"os"
	"path"
	"sync"

	"github.com/chaitin/libveinmind/go/plugin/log"
	_ "github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious/database/model"
)

var instance *gorm.DB
var once sync.Once

func GetDbInstance() *gorm.DB {
	once.Do(func() {
		databasePath := os.Getenv("DATABASE_PATH")
		wd, _ := os.Getwd()
		databasePath = path.Join(wd, databasePath)
		databasePathDir := path.Dir(databasePath)

		// 判断文件夹是否存在
		if _, err := os.Stat(databasePathDir); os.IsNotExist(err) {
			err := os.Mkdir(databasePathDir, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
		ist, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}

		instance = ist
	})

	return instance
}

func init() {
	Migrate()
}

func Migrate() {
	err := GetDbInstance().AutoMigrate(model.MaliciousFileInfo{}, model.ReportData{}, model.ReportImage{}, model.ReportLayer{})
	if err != nil {
		log.Error(err)
	}
}
