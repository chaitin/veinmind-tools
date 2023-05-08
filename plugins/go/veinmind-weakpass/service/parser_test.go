package service

import (
	"os"
	"testing"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-weakpass/model"
	"github.com/stretchr/testify/assert"
)

func TestTomcatParse(t *testing.T) {
	tomcat := &tomcatService{}
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
	redis := &redisService{}
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

func TestSshParse(t *testing.T) {
	Ssh := &SshService{}
	sshFile, err := os.Open("../test/shadow")
	assert.Nil(t, err)

	expectRecords := []model.Record{}
	expectRecords = append(expectRecords, model.Record{Username: "test",
		Password:   "$6$3oF7bkISmfCcnGIC$X588PbRFjkh5WDQfXcrLLYnYPN7bsjntaytebGGh3nsXp6d4uJCp3MCu54JSVXoZ8NxGWS5FxMcnloKvM4FXV/",
		Attributes: nil})
	expectRecords = append(expectRecords, model.Record{Username: "redis",
		Password:   "*",
		Attributes: nil})
	records, err := Ssh.GetRecords(sshFile)
	assert.Nil(t, err)
	assert.Equal(t, expectRecords, records)

}

func TestMyISAMParse(t *testing.T) {
	mysql := &mysql5Service{}

	mysqlMyd, err := os.Open("../test/user.MYD")
	assert.Nil(t, err)

	expectRecords := []model.Record{}
	expectRecords = append(expectRecords, model.Record{
		Username:   "root",
		Password:   "81f5e21e35407d884a6cd4a731aebfb6af209e1b",
		Attributes: nil,
	})

	records, err := mysql.GetRecords(mysqlMyd)
	assert.Nil(t, err)
	assert.Equal(t, len(expectRecords), len(records))

	for i, item := range records {
		assert.Nil(t, item.Attributes)
		assert.Equal(t, expectRecords[i].Username, item.Username)
		assert.Contains(t, item.Password, expectRecords[i].Password)
	}
}

func TestInnoDBParse(t *testing.T) {
	mysql := &mysql8Service{}

	mysqlIbd, err := os.Open("../test/mysql.ibd")
	assert.Nil(t, err)

	var expectRecords []model.Record
	expectRecords = append(expectRecords, model.Record{
		Username:   "root",
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

func TestVsftpdParse(t *testing.T) {
	vsftpd := &vsftpdService{}
	vsftpdFile, err := os.Open("../test/virtual_users.db")
	assert.Nil(t, err)

	expectRecords := []model.Record{}
	expectRecords = append(expectRecords, model.Record{
		Username:   "myuser",
		Password:   "mypass",
		Attributes: nil,
	})

	records, err := vsftpd.GetRecords(vsftpdFile)
	assert.Nil(t, err)
	assert.Equal(t, len(expectRecords), len(records))

	for i, item := range records {
		assert.Nil(t, item.Attributes)
		assert.Equal(t, expectRecords[i].Username, item.Username)
		assert.Contains(t, item.Password, expectRecords[i].Password)
	}
}

func TestProftpdParse(t *testing.T) {
	vsftpd := &proftpdService{}
	vsftpdFile, err := os.Open("../test/ftpd.passwd")
	assert.Nil(t, err)

	expectRecords := []model.Record{}
	expectRecords = append(expectRecords, model.Record{
		Username:   "user",
		Password:   "$1$U2Y3FMHr$NMXF3I.9Ym.lXkBBwGhLl",
		Attributes: nil,
	})

	records, err := vsftpd.GetRecords(vsftpdFile)
	assert.Nil(t, err)
	assert.Equal(t, len(expectRecords), len(records))

	for i, item := range records {
		assert.Nil(t, item.Attributes)
		assert.Equal(t, expectRecords[i].Username, item.Username)
		assert.Contains(t, item.Password, expectRecords[i].Password)
	}
}
