package dnsconfig

import (
	"log"
	"regexp"
	"strings"
)

func matchWildcard(wc, target string) bool {
	r := strings.Replace(wc, ".", "\\.", -1)
	r = strings.Replace(r, "+", "\\+", -1)
	r = strings.Replace(r, "*", "[^\\.]+", -1)
	re, err := regexp.Compile("^" + r + "$")
	if err != nil {
		log.Println("Could not make regexp from", wc, err)
		return false
	}
	if re.MatchString(target) {
		return true
	}
	return false
}
