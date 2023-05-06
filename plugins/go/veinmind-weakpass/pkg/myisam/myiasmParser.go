package myisam

import (
	"io"
)

const (
	// EmptyPasswordPlaceholder 用于没有修改过密码的 user page
	EmptyPasswordPlaceholder = "THISISNOTAVALIDPASSWORDTHATCANBEUSEDHERE"

	// LocalHost 用于识别Host是否为仅限本地登陆
	LocalHost = "localhost"
)

type MysqlInfo struct {
	Host     string
	Name     string
	Plugin   string
	Password string
}

// ParseUserFile 从文件中解析用户名和密码
func ParseUserFile(f io.Reader) (infos []MysqlInfo, err error) {
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	contentLen := len(content)
	idx := 0
	for idx < contentLen {
		record := dispatchRecord(content, idx)
		recType := record["rec_type"].(uint8)
		if 0 < recType && recType <= 6 {
			info := MysqlInfo{}
			res := parseRecord(content, record)
			info.Host = res["host"].(string)
			if res["native"] == true {
				info.Plugin = "mysql_native_password"
			} else {
				info.Plugin = "caching_sha2_password"
			}
			info.Name = res["user"].(string)
			info.Password = res["password"].(string)
			if info.Password != EmptyPasswordPlaceholder && info.Host != LocalHost {
				infos = append(infos, info)
			}
		}
		idx += record["block_len"].(int)
	}

	return infos, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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

func readRecord(content []byte, idx int, headerLen int, dataPos int, dataLen int, nextPos int, unusedLen int) (record map[string]interface{}) {
	recType := content[idx]
	dataLenValue := readLen(content, idx+dataPos, dataLen)
	unusedLenValue := unusedLen
	if unusedLen > 0 {
		unusedLenValue = int(content[idx+unusedLen])
	}
	blockLen := pad(headerLen + dataLenValue + unusedLenValue)
	nextRec := make(map[string]interface{})
	if nextPos > 0 {
		nextRec = dispatchRecord(content, readLen(content, idx+nextPos, 8))
	}
	record = map[string]interface{}{
		"rec_type":   recType,
		"block_len":  blockLen,
		"data_len":   dataLenValue,
		"next_rec":   nextRec,
		"data_begin": idx + headerLen,
	}
	return
}

func dispatchRecord(content []byte, idx int) (record map[string]interface{}) {
	recType := content[idx]
	switch recType {
	case 0:
		record = readRecord(content, idx, 20, 1, 3, -1, 0)
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
		record = readRecord(content, idx, 16, 5, 3, 9, 0)
	default:
		record = make(map[string]interface{})
	}
	return
}

func parseRecord(content []byte, rec map[string]interface{}) (result map[string]interface{}) {
	first := rec["data_begin"].(int) + 3
	hostLen := int(content[first])
	host := string(content[first+1 : first+1+hostLen])

	userLen := int(content[first+hostLen+1])
	user := string(content[first+hostLen+1+1 : first+hostLen+1+1+userLen])

	native := false
	passwordMaxLen := 40
	var password []byte
	idx := first + hostLen + 1 + 1 + userLen
	for {
		last := rec["data_begin"].(int) + rec["data_len"].(int)
		passwordLen := len(password)
		if passwordLen == 0 {
			for idx < last {
				if content[idx] == 21 {
					native = true
				}
				if content[idx] == 42 {
					break
				}
				idx++
			}
		}

		if idx+1 <= min(last, passwordMaxLen-passwordLen+idx+1) {
			password = append(password, content[idx+1:min(last, passwordMaxLen-passwordLen+idx+1)]...)
		}

		if len(rec["next_rec"].(map[string]interface{})) == 0 {
			break
		} else {
			rec = rec["next_rec"].(map[string]interface{})
			idx = rec["data_begin"].(int)
		}
	}
	result = map[string]interface{}{
		"host":     host,
		"user":     user,
		"password": string(password),
		"native":   native,
	}
	return
}
