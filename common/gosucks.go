package common

// a package because go can be freaking unfriendly at times

// returns true if the slice contains the string
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
