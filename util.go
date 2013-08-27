package main

import (
	"strings"
)

// TODO: remove code duplication
func rmAnsiEsc(s string) string {
	var l int
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' &&
			i+1 < len(s) &&
			s[i+1] == '[' {
			in := 2
			for {
				r := s[i+in]
				if r >= '0' && r <= '9' || r == ';' {
					in++
					continue
				}
				if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z') {
					in = 0
				}
				break
			}
			if in != 0 {
				i += in
				continue
			}
		}
		b[l] = s[i]
		l++
	}
	return string(b[:l])
}

func hasPrefixedSGI(s string) (l int) {
	if len(s) < 3 || s[0] != '\x1b' || s[1] != '[' {
		return 0
	}

	l = 2
	for l < len(s) {
		r := s[l]
		if r >= '0' && r <= '9' || r == ';' {
			l++
			continue
		}
		if r == 'm' {
			return l + 1
		}
		break
	}
	return 0
}

func lastEsc(s string) string {
	for {
		i := strings.LastIndex(s, "\x1b[")
		if i == -1 {
			return ""
		}
		l := hasPrefixedSGI(s[i:])
		if l != 0 {
			return s[i : i+l]
		}
		s = s[:i]
	}
}
