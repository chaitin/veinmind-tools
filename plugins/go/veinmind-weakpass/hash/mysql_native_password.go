package hash

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var _ Hash = (*MySQL)(nil)

type MySQL struct {
}

func (m *MySQL) ID() string {
	return "mysql"
}

func (m *MySQL) Match(hash, guess string) (flag bool, err error) {
	var checker Hash
	if len(hash) > 3 && hash[:3] == "$A$" {
		checker = &CachingSha2Password{}
	} else {
		checker = &MysqlNative{}
	}
	return checker.Match(hash, guess)
}

var _ Hash = (*MysqlNative)(nil)

type MysqlNative struct {
}

func (i *MysqlNative) ID() string {
	return "mysql_native_password"
}

func (i *MysqlNative) Match(hash, guess string) (flag bool, err error) {
	if strings.Contains(hash, "*") {
		r := sha1.Sum([]byte(guess))
		r = sha1.Sum(r[:])
		s := fmt.Sprintf("%x", r)
		if strings.Contains(hash, s) {
			return true, nil
		}
	}
	return false, errors.New("mysql_passwd: malformed entry ")
}

var _ Hash = (*CachingSha2Password)(nil)

type CachingSha2Password struct {
}

func (i *CachingSha2Password) ID() string {
	return "caching_sha2_password"
}

func (i *CachingSha2Password) Match(hash, guess string) (flag bool, err error) {
	if hash[:3] == "$A$" {
		return checkHashingPassword([]byte(hash[:70]), guess)
	}
	return false, errors.New("mysql_passwd: malformed entry ")
}

const (
	// MIXCHARS is the number of characters to use in the mix
	MIXCHARS = 32
	// SALT_LENGTH is the length of the salt
	SALT_LENGTH = 20 //nolint: revive
	// ITERATION_MULTIPLIER is the number of iterations to use
	ITERATION_MULTIPLIER = 1000 //nolint: revive
)

func b64From24bit(b []byte, n int, buf *bytes.Buffer) {
	b64t := []byte("./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	w := (int64(b[0]) << 16) | (int64(b[1]) << 8) | int64(b[2])
	for n > 0 {
		n--
		buf.WriteByte(b64t[w&0x3f])
		w >>= 6
	}
}

// sha256Hash is an util function to calculate sha256 hash.
func sha256Hash(input []byte) []byte {
	res := sha256.Sum256(input)
	return res[:]
}

// 'hash' function should return an array with 32 bytes, the same as SHA-256
// From https://github.com/pingcap/tidb/blob/ff78940594893cb970df58000dbc69cd5631e696/parser/auth/caching_sha2.go#L72 with some fix
func hashCrypt(plaintext string, salt []byte, iterations int, hash func([]byte) []byte) string {
	// Numbers in the comments refer to the description of the algorithm on https://www.akkadia.org/drepper/SHA-crypt.txt

	// 1, 2, 3
	bufA := bytes.NewBuffer(make([]byte, 0, 4096))
	bufA.Write([]byte(plaintext))
	bufA.Write(salt)

	// 4, 5, 6, 7, 8
	bufB := bytes.NewBuffer(make([]byte, 0, 4096))
	bufB.Write([]byte(plaintext))
	bufB.Write(salt)
	bufB.Write([]byte(plaintext))
	sumB := hash(bufB.Bytes())
	bufB.Reset()

	// 9, 10
	var i int
	for i = len(plaintext); i > MIXCHARS; i -= MIXCHARS {
		bufA.Write(sumB[:MIXCHARS])
	}
	bufA.Write(sumB[:i])

	// 11
	for i = len(plaintext); i > 0; i >>= 1 {
		if i%2 == 0 {
			bufA.Write([]byte(plaintext))
		} else {
			bufA.Write(sumB[:])
		}
	}

	// 12
	sumA := hash(bufA.Bytes())
	bufA.Reset()

	// 13, 14, 15
	bufDP := bufA
	for range []byte(plaintext) {
		bufDP.Write([]byte(plaintext))
	}
	sumDP := hash(bufDP.Bytes())
	bufDP.Reset()

	// 16
	p := make([]byte, 0, sha256.Size)
	for i = len(plaintext); i > 0; i -= MIXCHARS {
		if i > MIXCHARS {
			p = append(p, sumDP[:]...)
		} else {
			p = append(p, sumDP[0:i]...)
		}
	}

	// 17, 18, 19
	bufDS := bufA
	for i = 0; i < 16+int(sumA[0]); i++ {
		bufDS.Write(salt)
	}
	sumDS := hash(bufDS.Bytes())
	bufDS.Reset()

	// 20
	s := make([]byte, 0, 32)
	for i = len(salt); i > 0; i -= MIXCHARS {
		if i > MIXCHARS {
			s = append(s, sumDS[:]...)
		} else {
			s = append(s, sumDS[0:i]...)
		}
	}

	// 21
	bufC := bufA
	var sumC []byte
	for i = 0; i < iterations; i++ {
		bufC.Reset()
		if i&1 != 0 {
			bufC.Write(p)
		} else {
			bufC.Write(sumA[:])
		}
		if i%3 != 0 {
			bufC.Write(s)
		}
		if i%7 != 0 {
			bufC.Write(p)
		}
		if i&1 != 0 {
			bufC.Write(sumA[:])
		} else {
			bufC.Write(p)
		}
		sumC = hash(bufC.Bytes())
		sumA = sumC
	}

	// 22
	buf := bytes.NewBuffer(make([]byte, 0, 100))
	buf.Write([]byte{'$', 'A', '$'})
	rounds := fmt.Sprintf("%03X", iterations/ITERATION_MULTIPLIER)
	buf.Write([]byte(rounds))
	buf.Write([]byte{'$'})
	buf.Write(salt)

	b64From24bit([]byte{sumC[0], sumC[10], sumC[20]}, 4, buf)
	b64From24bit([]byte{sumC[21], sumC[1], sumC[11]}, 4, buf)
	b64From24bit([]byte{sumC[12], sumC[22], sumC[2]}, 4, buf)
	b64From24bit([]byte{sumC[3], sumC[13], sumC[23]}, 4, buf)
	b64From24bit([]byte{sumC[24], sumC[4], sumC[14]}, 4, buf)
	b64From24bit([]byte{sumC[15], sumC[25], sumC[5]}, 4, buf)
	b64From24bit([]byte{sumC[6], sumC[16], sumC[26]}, 4, buf)
	b64From24bit([]byte{sumC[27], sumC[7], sumC[17]}, 4, buf)
	b64From24bit([]byte{sumC[18], sumC[28], sumC[8]}, 4, buf)
	b64From24bit([]byte{sumC[9], sumC[19], sumC[29]}, 4, buf)
	b64From24bit([]byte{0, sumC[31], sumC[30]}, 3, buf)

	return buf.String()
}

// checkHashingPassword checks if a caching_sha2_password or tidb_sm3_password authentication string matches a password
// From https://github.com/pingcap/tidb/blob/ff78940594893cb970df58000dbc69cd5631e696/parser/auth/caching_sha2.go#L223 with some fix
func checkHashingPassword(pwHash []byte, password string) (bool, error) {
	pwHashParts := bytes.Split(pwHash, []byte("$"))
	if len(pwHashParts) < 3 {
		return false, errors.New("failed to decode hash parts")
	}

	hashType := string(pwHashParts[1])
	if hashType != "A" {
		return false, errors.New("digest type is incompatible")
	}

	iterations, err := strconv.ParseInt(string(pwHashParts[2]), 16, 64)
	if err != nil {
		return false, errors.New("failed to decode iterations")
	}
	iterations = iterations * ITERATION_MULTIPLIER
	salt := pwHashParts[3][:SALT_LENGTH]

	newHash := hashCrypt(password, salt, int(iterations), sha256Hash)

	return bytes.Equal(pwHash, []byte(newHash)), nil
}
