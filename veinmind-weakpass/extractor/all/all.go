package all

import (
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/extractor"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/extractor/mysql"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/extractor/redis"
	"github.com/chaitin/veinmind-tools/veinmind-weakpass/extractor/tomcat"
)

var All = []extractor.Extractor{
	&tomcat.Tomcat{},
	&redis.Redis{},
	&mysql.Mysql{},
}
