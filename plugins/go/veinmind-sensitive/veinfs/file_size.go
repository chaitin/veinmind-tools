package veinfs

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Size uint64

const (
	Bit  Size = 1
	Byte      = Bit << 3
	KB        = Byte << 10
	MB        = KB << 10
	GB        = MB << 10
	TB        = GB << 10
)

var (
	unitMap = map[string]Size{
		"b":  Bit,
		"B":  Byte,
		"KB": KB,
		"MB": MB,
		"GB": GB,
	}
	sizeMatchRegx = regexp.MustCompile(`^(\d+|\d+\.[\d+]{1,2})\s*(b|B|KB|MB|GB)$`)
	sizeParseRegx = regexp.MustCompile(`^\d+\.[\d+]{1,2}|\d+`)
	unitParseRegx = regexp.MustCompile(`(b|B|KB|MB|GB)$`)
)

func NewSize(text string) (Size, error) {
	if text == "" || text == "0" {
		return 0, nil
	}

	if !sizeMatchRegx.MatchString(text) {
		return 0, errors.New("not supported text")
	}

	unit := strings.TrimSpace(unitParseRegx.FindString(text))
	u, ok := unitMap[unit]
	if !ok {
		return 0, errors.New("not supported text")
	}

	size, err := strconv.ParseFloat(sizeParseRegx.FindString(text), 64)
	if err != nil {
		return 0, err
	}

	s := Size(size * float64(u))
	return s, nil
}

func (s Size) String() string {
	n := uint64(s)
	t := uint64(0b1111111111)
	us := []string{"PB", "TB", "GB", "MB", "KB", "Byte", "Bit"}
	vs := []uint64{
		(n & (t << 53)) >> 53, (n & (t << 43)) >> 43, (n & (t << 33)) >> 33, (n & (t << 23)) >> 23,
		(n & (t << 13)) >> 13, (n & (t << 03)) >> 3, n & 0b111,
	}

	for i := 0; i < len(us); i++ {
		if vs[i] == 0 && i != len(us)-1 {
			continue
		}

		if i == len(us)-1 || i == len(us)-2 {
			return fmt.Sprintf("%d %s", vs[i], us[i])
		} else {
			if vs[i+1] != 0 {
				return fmt.Sprintf("%d.%.0f %s", vs[i], (float64(vs[i+1])/1024)*10, us[i])
			} else {
				return fmt.Sprintf("%d %s", vs[i], us[i])
			}
		}
	}

	return ""
}
