module github.com/chaitin/veinmind-tools/veinmind-malicious

go 1.16

replace github.com/chaitin/veinmind-tools/veinmind-common/go v1.0.0 => ../veinmind-common/go

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20211005130812-5bb3c17173e5
	github.com/VirusTotal/vt-go v0.0.0-20211209151516-855a1e790678
	github.com/chaitin/libveinmind v1.0.7
	github.com/chaitin/veinmind-tools/veinmind-common/go v1.0.0
	github.com/dutchcoders/go-clamd v0.0.0-20170520113014-b970184f4d9e
	github.com/joho/godotenv v1.4.0
	github.com/mattn/go-sqlite3 v1.14.10 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.3.0
	golang.org/x/net v0.0.0-20220114011407-0dd24b26b47d // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	gorm.io/driver/sqlite v1.2.6
	gorm.io/gorm v1.22.5
)
