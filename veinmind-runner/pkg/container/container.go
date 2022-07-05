package container

import "os"

func InContainer() bool {
	if _, err := os.Open("/.dockerenv"); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
