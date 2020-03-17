package app

import (
	"strconv"
	"strings"
	"time"
)

func beUint(s string) uint {
	if len(s) == 0 {
		return 0
	}
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return uint(n)
}

func beTrace(n uint, h string) string {
	if n > 0 && len(h) >= 7 {
		return strconv.Itoa(int(n)) + "." + h[:7]
	}
	return ""
}

func beTime(s string) time.Time {
	if len(s) == 0 {
		return time.Time{}
	}

	var sec, nsec int64
	var err error
	// sec[.nsec] <=> 1572017404[.638981238]
	ss := strings.SplitN(s, ".", 2)
	if len(ss[0]) > 0 {
		sec, err = strconv.ParseInt(ss[0], 10, 64)
		if err != nil {
			panic(err)
		}
	}
	if len(ss) > 1 && len(ss[1]) > 0 {
		nsec, _ = strconv.ParseInt(ss[1], 10, 64)
	}

	return time.Unix(sec, nsec)
}
