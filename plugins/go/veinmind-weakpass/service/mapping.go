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
	case *mysql5Service:
		return event.Mysql
	case *mysql8Service:
		return event.Mysql
	case *tomcatService:
		return event.Tomcat
	case *vsftpdService:
		return event.FTP
	case *proftpdService:
		return event.FTP
	default:
		return 0
	}
}
