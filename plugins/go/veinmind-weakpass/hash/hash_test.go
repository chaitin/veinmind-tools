package hash

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashMysqlNative(t *testing.T) {
	mysqlNative := &MysqlNative{}
	// password which got from file mysql.ibd in docker image
	stringGotFromFile := "*6bb4837eb74329105ee4568dda7dc67ed2ca2ad9"
	stringInDict := "123456"
	find, _ := mysqlNative.Match(stringGotFromFile, stringInDict)
	assert.True(t, find)

	stringGotFromFile = "6bb4837eb74329105ee4568dda7dc67ed2ca2ad9"
	find, _ = mysqlNative.Match(stringGotFromFile, stringInDict)
	assert.False(t, find)
}

func TestHashCachingSha2Password(t *testing.T) {
	cachingSha2Password := &CachingSha2Password{}
	// password which got from file mysql.ibd in docker image
	b64PwdFromFile := "JEEkMDA1JFodUGBzOkJlVwx/bwRYfwg8f3hoV3ZxRUN3c3JGa3pwa0kuVWxXbEpSeEM1ZkVYVUhscU42WFVuaGpHTS5zNwFixZPBAQICRkAAAAAgfAAAOPokJQ=="
	stringInDict := "123456"
	pwd, err := base64.StdEncoding.DecodeString(b64PwdFromFile)
	if err != nil {
		t.Fatal(err)
	}
	find, _ := cachingSha2Password.Match(string(pwd), stringInDict)
	assert.True(t, find)
	find, _ = cachingSha2Password.Match(string(pwd), "12345678")
	assert.False(t, find)
}

func TestHashShadow(t *testing.T) {
	shadow := &Shadow{}
	stringGotFromShadow := "$6$3oF7bkISmfCcnGIC$X588PbRFjkh5WDQfXcrLLYnYPN7bsjntaytebGGh3nsXp6d4uJCp3MCu54JSVXoZ8NxGWS5FxMcnloKvM4FXV/"
	stringInDict := "123456"
	find, _ := shadow.Match(stringGotFromShadow, stringInDict)
	assert.True(t, find)
	stringGotFromShadow = "$5$3oF7bkISmfCcnGIC$X588PbRFjkh5WDQfXcrLLYnYPN7bsjntaytebGGh3nsXp6d4uJCp3MCu54JSVXoZ8NxGWS5FxMcnloKvM4FXV/"
	find, _ = shadow.Match(stringGotFromShadow, stringInDict)
	assert.False(t, find)
	stringGotFromShadow = "$y$j9T$u.6g7OsvLMeIhyBve4MsT/$S8btad/SN91QUwrQjf84XEzLj2w64xVA1UQQN6xZeQ2"
	find, _ = shadow.Match(stringGotFromShadow, stringInDict)
	assert.True(t, find)
	stringGotFromShadow = "$y$j9T$u.6g01svLMeIhyBve4MsT/$S8btad/SN91QUwrQjf84XEzLj2w64xVA1UQQN6xZeQ2"
	find, _ = shadow.Match(stringGotFromShadow, stringInDict)
	assert.False(t, find)

}

func TestHashPlain(t *testing.T) {
	plain := &Plain{}
	stringGotFromFile := "admin"
	stringInDict := "admin"
	find, _ := plain.Match(stringGotFromFile, stringInDict)
	assert.True(t, find)
	stringGotFromFile = "tomcat_admin"
	find, _ = plain.Match(stringGotFromFile, stringInDict)
	assert.False(t, find)
}
