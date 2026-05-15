package dnsconfig

import (
	"log"
	"regexp"
	"strings"
	"sync"
)

var regexCache sync.Map

func matchWildcard(wc, target string) bool {
	cached, ok := regexCache.Load(wc)
	if !ok {
		re, err := compileWildcard(wc)
		if err != nil {
			log.Println("Could not make regexp from", wc, err)
			return false
		}
		cached, _ = regexCache.LoadOrStore(wc, re)
	}
	return cached.(*regexp.Regexp).MatchString(target)
}

func compileWildcard(wc string) (*regexp.Regexp, error) {
	r := wc
	if !strings.HasPrefix(wc, "^") && !strings.HasSuffix(wc, "$") {
		parts := strings.Split(wc, "*")
		for i, p := range parts {
			parts[i] = regexp.QuoteMeta(p)
		}
		r = "^" + strings.Join(parts, `[^.]+`) + "$"
	}
	return regexp.Compile(r)
}
