package database

import (
	_ "github.com/chaitin/veinmind-tools/veinmind-malicious/config"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/database/model"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/sdk/common"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
	"sync"
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
				common.Log.Fatal(err)
			}
		}
		ist, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
		if err != nil {
			common.Log.Fatal(err)
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
		common.Log.Error(err)
	}
}
