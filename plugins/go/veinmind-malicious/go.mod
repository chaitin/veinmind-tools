module github.com/chaitin/veinmind-tools/plugins/go/veinmind-malicious

go 1.16

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20211005130812-5bb3c17173e5
	github.com/VirusTotal/vt-go v0.0.0-20211209151516-855a1e790678
	github.com/chaitin/libveinmind v1.3.1
	github.com/chaitin/veinmind-common-go v1.1.9
	github.com/joho/godotenv v1.4.0
	github.com/mattn/go-sqlite3 v1.14.10 // indirect
	github.com/spf13/cobra v1.5.0
	github.com/testwill/go-clamd v1.0.0
	gorm.io/driver/sqlite v1.2.6
	gorm.io/gorm v1.22.5
)

replace google.golang.org/grpc/naming => github.com/xiegeo/grpc-naming v1.29.1-alpha
