package myisam

import (
	"errors"
	"io"
)

const (
	// EmptyPasswordPlaceholder 用于没有修改过密码的 user page
	EmptyPasswordPlaceholder = "*THISISNOTAVALIDPASSWORDTHATCANBEUSEDHERE"

	// LocalHost 用于识别Host是否为仅限本地登陆
	LocalHost = "localhost"
)

type UserInfo struct {
	Host     string
	Name     string
	Plugin   string
	Password string
}

type Record struct {
	RecType   uint8
	BlockLen  int
	DataLen   int
	NextRec   *Record
	DataBegin int
}

// ParseUserFile 从文件中解析用户名和密码
func ParseUserFile(f io.Reader) (infos []*UserInfo, err error) {
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	contentLen := len(content)
	idx := 0
	for idx < contentLen {
		record := dispatchRecord(content, idx)
		recType := record.RecType
		if 0 < recType && recType <= 6 {
			if info, err := parseRecord(content, record); err != nil {
				return infos, err
			} else {
				infos = append(infos, info)
			}
		}
		idx += record.BlockLen
	}

	return infos, nil
}

func readLen(content []byte, begin int, l int) int {
	sumLen := 0
	for _, bit := range content[begin : begin+l] {
		sumLen = (sumLen << 8) + int(bit)
	}
	return sumLen
}

// 四位补齐一位字节
func pad(dataLen int) int {
	byteLen := dataLen >> 2
	return (byteLen + ((dataLen - (byteLen << 2)) & 1)) << 2
}

func readRecord(content []byte, idx int, headerLen int, dataPos int, dataLen int, nextPos int, unusedLen int) (record *Record) {
	recType := content[idx]
	dataLenValue := readLen(content, idx+dataPos, dataLen)
	unusedLenValue := unusedLen
	if unusedLen > 0 {
		unusedLenValue = int(content[idx+unusedLen])
	}
	blockLen := pad(headerLen + dataLenValue + unusedLenValue)
	var nextRec *Record
	if nextPos > 0 {
		nextRec = dispatchRecord(content, readLen(content, idx+nextPos, 8))
	}

	record = &Record{
		RecType:   recType,
		BlockLen:  blockLen,
		DataLen:   dataLenValue,
		NextRec:   nextRec,
		DataBegin: idx + headerLen,
	}

	return
}

// Related info: https://github.com/twitter-forks/mysql/blob/865aae5f23e2091e1316ca0e6c6651d57f786c76/storage/myisam/mi_dynrec.c#LL1890C1-L1890C1
func dispatchRecord(content []byte, idx int) (record *Record) {
	recType := content[idx]
	switch recType {
	case 0:
		record = readRecord(content, idx, 0, 1, 3, -1, 0)
	case 1:
		record = readRecord(content, idx, 3, 1, 2, -1, 0)
	case 2:
		record = readRecord(content, idx, 4, 1, 3, -1, 0)
	case 3:
		record = readRecord(content, idx, 4, 1, 2, -1, 3)
	case 4:
		record = readRecord(content, idx, 5, 1, 3, -1, 4)
	case 5:
		record = readRecord(content, idx, 13, 3, 2, 5, 0)
	case 6:
		record = readRecord(content, idx, 15, 4, 3, 7, 0)
	case 7:
		record = readRecord(content, idx, 3, 1, 2, -1, 0)
	case 8:
		record = readRecord(content, idx, 4, 1, 3, -1, 0)
	case 9:
		record = readRecord(content, idx, 4, 1, 2, -1, 3)
	case 10:
		record = readRecord(content, idx, 5, 1, 3, -1, 4)
	case 11:
		record = readRecord(content, idx, 11, 1, 2, 3, 0)
	case 12:
		record = readRecord(content, idx, 12, 1, 3, 4, 0)
	case 13:
		record = readRecord(content, idx, 16, 5, 3, 8, 0)
	default:
		record = nil
	}
	return
}

var RecordDataBrokenErr = errors.New("record data broken")

// parseRecord parse user.MYD, return nil if error happened
func parseRecord(content []byte, rec *Record) (result *UserInfo, err error) {
	var trueContent []byte
	for rec != nil {
		// 这里应该根据rec还原数据 content的内容是不能直接使用的
		// Related info: https://github.com/twitter-forks/mysql/blob/865aae5f23e2091e1316ca0e6c6651d57f786c76/storage/myisam/mi_dynrec.c#LL1890C1-L1890C1
		trueContent = append(trueContent, content[rec.DataBegin:rec.DataBegin+rec.DataLen]...)
		rec = rec.NextRec
	}

	hostLenPos := 3
	if len(trueContent) <= hostLenPos {
		return nil, RecordDataBrokenErr
	}
	hostLen := int(trueContent[hostLenPos])

	userLenPos := hostLenPos + hostLen + 1 // 可以简化 但是这么写便于阅读
	if len(trueContent) <= userLenPos {
		return nil, RecordDataBrokenErr
	}
	userLen := int(trueContent[userLenPos])

	if len(trueContent) <= userLenPos+1+userLen {
		return nil, RecordDataBrokenErr
	}
	host := string(trueContent[hostLenPos+1 : hostLenPos+1+hostLen])
	user := string(trueContent[userLenPos+1 : userLenPos+1+userLen])

	// caching_sha2_password在>=5.6版本被支持 也就是自5.6版开始 user表结构有所变化，加入authentication_string列
	// authentication_string列有一个固定特征是长度为21 且跟在user列后面（user数据结束后会有很长一段2）
	// 如果是<5.6版本的表，user列后面直接紧跟着就是42了（*）
	plugin := "mysql_native_password"
	passwdLenPos := userLenPos + userLen + 1 // 先当作是<5.6的版本提取
	var password string
	if passwdLen := int(trueContent[passwdLenPos]); passwdLen == 42 { // 判断出来不是>=5.6的版本
		password = string(trueContent[passwdLenPos : passwdLenPos+41])
	} else { // >=5.6版本 authentication_string列在前 password列在后
		passwdLenPos++
		for {
			if passwdLenPos < len(trueContent) {
				if trueContent[passwdLenPos] == 21 { // 读到了authentication_string列
					plugin = string(trueContent[passwdLenPos+1 : passwdLenPos+1+21])
					passwdLenPos += 1 + 21
					passwdLen = int(trueContent[passwdLenPos]) // TEXT类型 这里认为长度密码字段长度不会超过256 不计算第二位 https://dev.mysql.com/doc/refman/8.0/en/storage-requirements.html
					password = string(trueContent[passwdLenPos+2 : passwdLenPos+2+passwdLen])
					break
				}
				passwdLenPos++
			} else {
				return nil, RecordDataBrokenErr
			}
		}
	}
	return &UserInfo{
		Host:     host,
		Name:     user,
		Password: password,
		Plugin:   plugin,
	}, nil
}
