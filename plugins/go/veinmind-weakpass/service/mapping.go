package service

import (
	"github.com/chaitin/veinmind-common-go/service/report/event"
)

func GetType(service IService) event.WeakpassService {
	switch service.(type) {
	case *SshService:
		return event.SSH
	case *redisService:
		return event.Redis
	case *mysqlService:
		return event.Mysql
	case *tomcatService:
		return event.Tomcat
	default:
		return 0
	}
}
