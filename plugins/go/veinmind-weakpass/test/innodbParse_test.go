package test

import (
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/pkg/innodb"
	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/pkg/myisam"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseMySQLInnoDB(t *testing.T) {
	f, err := os.Open("./mysql8.ibd")
	if err != nil {
		return
	}
	defer f.Close()

	page, err := innodb.FindUserPage(f)
	infos, err := innodb.ParseUserPage(page.Pagedata)
	if err != nil {
		return
	}
	for _, info := range infos {
		if info.Password[:3] == "$A$" {
			t.Log("caching_sha2_password")
		} else {
			t.Log("mysql_native_password")
		}
		t.Log(info)
	}
}

func TestParseMySQLInnoDB5(t *testing.T) {
	f, err := os.Open("./mysql.ibd")
	if err != nil {
		return
	}
	defer f.Close()

	page, err := innodb.FindUserPage(f)
	infos, err := innodb.ParseUserPage(page.Pagedata)
	if err != nil {
		return
	}
	for _, info := range infos {
		if info.Password[:3] == "$A$" {
			t.Log("caching_sha2_password")
		} else {
			t.Log("mysql_native_password")
		}
		t.Log(info)
	}
}

func TestParseMyISAM55(t *testing.T) {
	mysqlMyd, err := os.Open("../test/mysql5_5_myisam.MYD")
	assert.Nil(t, err)

	infos, err := myisam.ParseUserFile(mysqlMyd)
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, infos, 2)
	assert.Equal(t, "*6BB4837EB74329105EE4568DDA7DC67ED2CA2AD9", infos[0].Password)
}

func TestParseMyISAM56(t *testing.T) {
	mysqlMyd, err := os.Open("../test/user.MYD")
	assert.Nil(t, err)

	infos, err := myisam.ParseUserFile(mysqlMyd)
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, infos, 4)
	assert.Equal(t, "*81F5E21E35407D884A6CD4A731AEBFB6AF209E1B", infos[0].Password)
}
