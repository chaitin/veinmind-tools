package vregex

import (
	"regexp"
	"sync"
)

var (
	regexMap = sync.Map{}
)

// getRegexp returns *regexp.Regexp object with given <pattern>.
// It uses cache to enhance the performance for compiling regular expression pattern,
// which means, it will return the same *regexp.Regexp object with the same regular
// expression pattern.
//
// It is concurrent-safe for multiple goroutines.
func getRegexp(pattern string) (regex *regexp.Regexp, err error) {
	// Retrieve the regular expression object using reading lock.
	loaded, ok := regexMap.Load(pattern)
	if ok {
		return loaded.(*regexp.Regexp), nil
	} else {
		// If it does not exist in the cache,
		// it compiles the pattern and creates one.
		regex, err = regexp.Compile(pattern)
		if err != nil {
			return
		}
		// Cache the result object using writing lock.
		regexMap.Store(pattern, regex)
		return
	}
}
