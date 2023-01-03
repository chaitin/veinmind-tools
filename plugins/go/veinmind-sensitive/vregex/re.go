package vregex

// IsMatch checks whether given bytes <src> matches <pattern>.
func IsMatch(pattern string, src []byte) bool {
	if r, err := getRegexp(pattern); err == nil {
		return r.Match(src)
	}
	return false
}

// IsMatchString checks whether given string <src> matches <pattern>.
func IsMatchString(pattern string, src string) bool {
	return IsMatch(pattern, []byte(src))
}

// FindIndex find given bytes <src> index.
func FindIndex(pattern string, src []byte) []int {
	if r, err := getRegexp(pattern); err == nil {
		return r.FindIndex(src)
	}
	return nil
}

// FindStringIndex find given bytes <src> index.
func FindStringIndex(pattern string, src string) []int {
	if r, err := getRegexp(pattern); err == nil {
		return r.FindStringIndex(src)
	}
	return nil
}

// FindIndexWithContextContent find given bytes <src> index with context content.
func FindIndexWithContextContent(pattern string, src []byte, size int) (content []byte, loc []int) {
	rang := FindIndex(pattern, src)
	contextRange := make([]int, 2)
	contextHighlightRange := make([]int, 2)
	if rang != nil {
		if rang[0]-size <= 0 {
			contextRange[0] = 0
			contextHighlightRange[0] = rang[0]
		} else {
			contextRange[0] = rang[0] - size
			contextHighlightRange[0] = size
		}

		if rang[1]+size >= len(src) {
			contextRange[1] = len(src)
			contextHighlightRange[1] = contextHighlightRange[0] + rang[1] - rang[0]
		} else {
			contextRange[1] = rang[1] + size
			contextHighlightRange[1] = contextHighlightRange[0] + rang[1] - rang[0]
		}
		content = src[contextRange[0]:contextRange[1]]
		loc = contextHighlightRange
		return
	}

	return nil, nil
}
