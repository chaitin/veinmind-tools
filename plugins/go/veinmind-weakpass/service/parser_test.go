package service

import (
	"os"
	"testing"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
	"github.com/stretchr/testify/assert"
)

func TestTomcatParse(t *testing.T) {
	tomcat, err := GetModuleByName("tomcat")
	assert.Nil(t, err)

	tomcatfile, err := os.Open("../test/tomcat-users.xml")
	assert.Nil(t, err)

	expectRecords := []model.Record{}
	expectRecords = append(expectRecords, model.Record{
		Username: "both", Password: "tomcat",
		Attributes: map[string]string{"roles": "tomcat,role1"}})
	expectRecords = append(expectRecords, model.Record{
		Username: "role1", Password: "tomcat",
		Attributes: map[string]string{"roles": "role1"}})
	records, err := tomcat.GetRecords(tomcatfile)
	assert.Nil(t, err)
	assert.Equal(t, expectRecords, records)

}

func TestRedisParse(t *testing.T) {
	redis, err := GetModuleByName("redis")
	assert.Nil(t, err)

	redisFile, err := os.Open("../test/redis.conf")
	assert.Nil(t, err)

	expectRecords := []model.Record{}
	expectRecords = append(expectRecords, model.Record{
		Username: "", Password: "123456",
		Attributes: nil,
	})
	expectRecords = append(expectRecords, model.Record{
		Username: "", Password: "foobared",
		Attributes: nil,
	})
	records, err := redis.GetRecords(redisFile)
	assert.Nil(t, err)

	assert.Equal(t, expectRecords, records)

}

func TestShadowParse(t *testing.T) {
	shadow, err := GetModuleByName("ssh")
	assert.Nil(t, err)

	shadowFile, err := os.Open("../test/shadow")
	assert.Nil(t, err)

	expectRecords := []model.Record{}
	expectRecords = append(expectRecords, model.Record{Username: "test",
		Password:   "$6$3oF7bkISmfCcnGIC$X588PbRFjkh5WDQfXcrLLYnYPN7bsjntaytebGGh3nsXp6d4uJCp3MCu54JSVXoZ8NxGWS5FxMcnloKvM4FXV/",
		Attributes: nil})
	expectRecords = append(expectRecords, model.Record{Username: "redis",
		Password:   "*",
		Attributes: nil})
	records, err := shadow.GetRecords(shadowFile)
	assert.Nil(t, err)
	assert.Equal(t, expectRecords, records)

}

func TestMysqlParse(t *testing.T) {
	mysql, err := GetModuleByName("mysql")
	assert.Nil(t, err)

	mysqlIbd, err := os.Open("../test/mysql.ibd")
	assert.Nil(t, err)

	expectRecords := []model.Record{}
	expectRecords = append(expectRecords, model.Record{
		Username:   "root@localhost",
		Password:   "*6a7a490fb9dc8c33c2b025a91737077a7e9cc5e5",
		Attributes: nil,
	})
	records, err := mysql.GetRecords(mysqlIbd)
	assert.Nil(t, err)
	assert.Equal(t, len(expectRecords), len(records))

	for i, item := range records {
		assert.Nil(t, item.Attributes)
		assert.Equal(t, expectRecords[i].Username, item.Username)
		assert.Contains(t, item.Password, expectRecords[i].Password)
	}

}
