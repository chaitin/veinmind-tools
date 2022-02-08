package database

import (
	"fmt"
	"github.com/chaitin/veinmind-tools/veinmind-malicious/database/model"
	"testing"
)

func TestGetDbInstance(t *testing.T) {
	r := model.ReportImage{}
	GetDbInstance().Where("image_id = ?", "sha256:2a06e9574854723b606ef03dda327158edd3525c8454661712db49a5548bc0d0").Find(&r)
	//GetDbInstance().First(&r)
	fmt.Println(r)
}
