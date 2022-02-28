package ssh_passwd

// #include <stdlib.h>
// #include <unistd.h>
// #include <string.h>
// #define __USE_GNU 1
// #include <crypt.h>
// #cgo LDFLAGS: -Wl,-Bstatic -lcrypt -Wl,-Bdynamic
//
// static int passwd_match
// (const char* salt, const char* src, const char* dst) {
//     struct crypt_data cd;
//     memset(&cd, 0, sizeof(cd));
//     char* result = crypt_r(src, salt, &cd);
//     if(result == NULL) return 0;
//     return strcmp(result, dst) == 0;
// }
import "C"

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"unsafe"
)

var ErrMalformed = errors.New("ssh_passwd: malformed entry")

// PasswordMethod represents the entryption method.
type PasswordMethod uint8

const (
	// Password stored inside shadow file. Can only appear
	// inside ssh_passwd.Passwd.
	Shadowed PasswordMethod = iota

	// User is locked and cannot login.
	Locked

	// User does not have a password.
	NoPassword

	// Password encrypted with MD5. (=1)
	MD5

	// Password encrypted with blowfish. (=2/2a)
	Blowfish

	// Password encrypted with SHA256. (=5)
	SHA256

	// Password encrypted with SHA512. (=6)
	SHA512

	// NT-Hash (actually MD4)
	NTHash
)

// Password represents the data field inside.
type Password struct {
	// Method of current password.
	Method PasswordMethod

	// MethodString of the current password.
	MethodString string

	// Salt inside the password.
	Salt string

	// Hash hased value of the password.
	Hash string
}

// regexpKey make sure the validity of the salt and hash.
var regexpKey = regexp.MustCompile("[A-Za-z0-9./]+")

// Parse the password from the password field.
func ParsePassword(pass *Password, phrase string) error {
	// Eliminate the case that the password is shadow
	// and locked.
	switch phrase {
	case "x": // Stored in shadow.
		pass.Method = Shadowed
		return nil
	case "*", "!", "!!": // Cannot login with password.
		pass.Method = Locked
		return nil
	case "":
		pass.Method = NoPassword
		return nil
	default: // Continue on parsing.
	}

	// Split the string into multiple portions.
	s := strings.Split(phrase, "$")
	if len(s) != 4 {
		// $method$salt$hash => ["", method, salt, hash]
		return ErrMalformed
	}
	method := s[1]
	salt := s[2]
	hash := s[3]

	// Method should be recognized by the application.
	switch method {
	case "1":
		pass.Method = MD5
	case "2", "2a", "2y":
		pass.Method = Blowfish
		if len(hash) <= 22 {
			return ErrMalformed
		}
		salt = fmt.Sprintf("%s$%s", salt, hash[:22])
		hash = hash[22:]
	case "5":
		pass.Method = SHA256
	case "6":
		pass.Method = SHA512
	default:
		return ErrMalformed
	}
	pass.MethodString = method

	// Make sure the pass phrases are only made up of regexKey.
	if !regexpKey.MatchString(salt) {
		return ErrMalformed
	}
	pass.Salt = salt
	if !regexpKey.MatchString(hash) {
		return ErrMalformed
	}
	pass.Hash = hash
	return nil
}

// Match tests whether one of your guesses matches one of the given
// password. It can only be used for weak password detection purpose.
func (pw *Password) Match(guesses []string) (string, bool) {
	// You must not attempt to match shadow or locked password.
	if pw.Method == Shadowed || pw.Method == Locked {
		return "", false
	}

	// Match only empty guesses. (Though we should only guess
	// password with empty key checking.
	if pw.Method == NoPassword {
		for _, guess := range guesses {
			if guess == "" {
				return "", true
			}
		}
		return "", false
	}

	salt := fmt.Sprintf("$%s$%s", pw.MethodString, pw.Salt)
	// Blowfish not in mainline glibc, added in some Linux distributions
	if pw.Method == Blowfish {
		hash := []byte(fmt.Sprintf("%s%s", salt, pw.Hash))
		for _, guess := range guesses {
			if bcrypt.CompareHashAndPassword(hash, []byte(guess)) == nil {
				return guess, true
			}
		}
		return "", false
	}
	// Attempt to encrypt the password.
	saltCString := C.CString(salt)
	defer C.free(unsafe.Pointer(saltCString))
	dst := fmt.Sprintf("%s$%s", salt, pw.Hash)
	dstCString := C.CString(dst)
	defer C.free(unsafe.Pointer(dstCString))
	for _, guess := range guesses {
		srcCString := C.CString(guess)
		i := C.passwd_match(saltCString, srcCString, dstCString)
		C.free(unsafe.Pointer(srcCString))
		if i == C.int(1) {
			return guess, true
		}
	}
	return "", false
}
