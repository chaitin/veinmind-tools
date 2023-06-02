package innodb

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	PageSize = 0x4000
	// EmptyPasswordPlaceholder 用于没有修改过密码的 user page
	EmptyPasswordPlaceholder = "THISISACOMBINATIONOFINVALIDSALTANDPASSWORDTHATMUSTNEVERBRBEUSED"
	// PluginNameNative 用于mysql 8.0.11 版本定位 user page
	PluginNameNative = "mysql_native_password"
	// PluginNameCaching 用于mysql 8.0.13 - 19 版本的 user page
	PluginNameCaching = "caching_sha2_password"
	MysqlSYS          = "mysql.sys"
	// HostLengthBefore29 mysql 8.0.29 之前的 host 字段的长度
	HostLengthBefore29 int16 = 60
	// HostLengthAfter29 mysql 8.0.11 - 19 版本的 host 字段的长度
	HostLengthAfter29 int16 = 255
	FileHeaderSize          = 0x38
	PageHeaderSize          = 0x58
	ListSize                = 0x26
	PageDataSize            = 0x3f52
	PageTailSize            = 0x8
)

type FileHeader struct {
	FIL_PAGE_SPACE_OR_CHKSUM         uint32
	FIL_PAGE_OFFSET                  uint32
	FIL_PAGE_PREV                    uint32
	FIL_PAGE_NEXT                    uint32
	FIL_PAGE_LSN                     uint64
	FIL_PAGE_TYPE                    uint16
	FIL_PAGE_FILE_FLUSH_LSN          uint64
	FIL_PAGE_ARCH_LOG_NO_OR_SPACE_ID uint32
}

type PageHeader struct {
	PAGE_N_DIR_SLOTS  uint16
	PAGE_HEAP_TOP     uint16
	PAGE_N_HEAP       uint16
	PAGE_FREE         uint16
	PAGE_GARBAGE      uint16
	PAGE_LAST_INSERT  uint16
	PAGE_DIRECTION    uint16
	PAGE_N_DIRECTION  uint16
	PAGE_N_RECS       uint16
	PAGE_MAX_TRX_ID   uint64
	PAGE_LEVEL        uint16
	PAGE_INDEX_ID     uint64
	PAGE_BTR_SEG_LEAF [10]byte
	PAGE_BTR_SEG_TOP  [10]byte
}

type RecordHeader struct {
	NeedtoParse [3]byte
	NextRecord  int16
}
type InfimumSupremumRecord struct {
	Recordheader RecordHeader
	Text         [8]byte
}
type PageData struct {
	Infimum_record  InfimumSupremumRecord
	Supremum_record InfimumSupremumRecord
	Content         [PageDataSize - 0x26]byte
}

type PageTail struct {
	Content [PageTailSize]byte
}
type Page struct {
	Fileheader FileHeader
	Pageheader PageHeader
	Pagedata   PageData
	Pagetail   PageTail
}

// IsUserPage 通过用户名和加密方式的组合判断是否是 user 页面
func IsUserPage(buf []byte) bool {
	if bytes.Contains(buf, []byte(MysqlSYS)) {
		if bytes.Contains(buf, []byte(PluginNameCaching)) || bytes.Contains(buf, []byte(PluginNameNative)) {
			return true
		}
	}
	return false
}

// FindUserPage 从 mysql.ibd 中定位到 user 表的页面,并返回 user.ibd
func FindUserPage(f io.Reader) (page Page, err error) {
	r := bufio.NewReader(f)
	buf := make([]byte, 0, PageSize)
	foundPage := false
	for {
		n, err := io.ReadFull(r, buf[:cap(buf)])
		buf = buf[:n]
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			if !errors.Is(err, io.ErrUnexpectedEOF) {
				return page, err
			}
		}

		if IsUserPage(buf) {
			reader := bytes.NewReader(buf)
			binary.Read(reader, binary.BigEndian, &page)
			foundPage = true
			break
		}
	}
	if !foundPage {
		return page, errors.New(fmt.Sprintf("%s was not found on any page", MysqlSYS))
	}

	return page, nil
}

func File2Page(file string) (page Page, err error) {
	buf, err := os.Open(file)
	if err != nil {
		return
	}
	defer buf.Close()
	if err != nil {
		return page, err
	}
	binary.Read(buf, binary.BigEndian, &page)
	return page, nil
}

func Bytes2Int16(b []byte) int16 {
	bytesBuffer := bytes.NewBuffer(b)
	var x int16
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

type MysqlInfo struct {
	Host     string
	Name     string
	Plugin   string
	Password string
}

// ParseUserPage 从 user.ibd 中提取 user、host、plugin 和 password []bytes
func ParseUserPage(pagedata PageData) (infos []MysqlInfo, err error) {
	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.BigEndian, &pagedata)
	data := buf.Bytes()
	if err != nil {
		return infos, err
	}
	InfimumPos := pagedata.Infimum_record.Recordheader.NextRecord
	// 最小记录起始位置, 0xC - 0x8
	var startPos int16 = 0x5
	var endPos int16 = 0x12
	nextRecord := startPos + InfimumPos
	// 根据链表指针获取每行的数据
	for i := nextRecord; i != endPos; {
		// 这里对于正常的innodb文件解析下标不存在越界问题
		// 为防止解析文件出错导致程序中断,加入首尾边界检查
		if i < 2 || i+60+246 > 0x3f52 {
			return infos, errors.New("user page parse error")
		}
		nextRecordBytes := data[i-2 : i]
		next := Bytes2Int16(nextRecordBytes)
		var hostLen int16 = 60
		// var hostLen int16 = 255
		info := MysqlInfo{}
		Host := strings.TrimSpace(string(data[i : i+hostLen]))
		Name := strings.TrimSpace(string(data[i+hostLen : i+hostLen+32]))
		Plugin := strings.TrimSpace(string(data[i+hostLen+91 : i+hostLen+155]))
		Password := strings.TrimSpace(string(data[i+hostLen+155 : i+hostLen+246]))
		// 对于mysql 8.0.29 版本, host 长度为 255
		if Name == "" || Plugin == "" {
			if i+255+246 > 0x3f52 {
				return infos, errors.New("user page parse error")
			}
			hostLen = 255
			Host = strings.TrimSpace(string(data[i : i+hostLen]))
			Name = strings.TrimSpace(string(data[i+hostLen : i+hostLen+32]))
			Plugin = strings.TrimSpace(string(data[i+hostLen+91 : i+hostLen+155]))
			Password = strings.TrimSpace(string(data[i+hostLen+155 : i+hostLen+246]))
		}
		info.Name = Name
		info.Host = Host
		info.Plugin = Plugin
		info.Password = Password
		infos = append(infos, info)
		i += next
	}
	return infos, nil
}
