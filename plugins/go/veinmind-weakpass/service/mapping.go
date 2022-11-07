package service

import (
	"github.com/chaitin/veinmind-common-go/service/report"
)

func GetType(service IService) report.WeakpassService {
	switch service.(type) {
	case *SshService:
		return report.SSH
	case *redisService:
		return report.Redis
	case *mysqlService:
		return report.Mysql
	case *tomcatService:
		return report.Tomcat
	default:
		return 0
	}
}
