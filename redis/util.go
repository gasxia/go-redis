package redis

import (
	"math"
)

/* Convert a string into a long long. Returns 1 if the string could be parsed
 * into a (non-overflowing) long long, 0 otherwise. The value will be set to
 * the parsed value when appropriate. */
func string2ll(s string) (value int64, flag bool) {
	p := s
	plen := 0
	slen := len(s)
	var v uint64
	if slen == 0 {
		return 0, false
	}

	if slen == 1 && p[0] == '0' {
		return 0, true
	}
	var negative bool
	if p[0] == '-' {
		negative = true
		plen++
		p = p[1:]
		if plen == slen {
			return 0, true
		}
	}

	if p[0] >= '1' && p[0] <= '9' {
		v = uint64(p[0] - '0')
		plen++
		p = p[1:]
	} else {
		return 0, false
	}

	var r uint64
	for plen < slen && p[0] >= '0' && p[0] <= '9' {
		if v > math.MaxUint64/10 {
			return 0, false
		}
		v *= 10

		r = uint64(p[0] - '0')
		if v > (math.MaxUint64 - r) {
			return 0, false
		}
		v += r
		plen++
		p = p[1:]
	}

	if plen < slen {
		return 0, false
	}
	if negative {
		if v > uint64(math.MaxInt64+1) { // Overflow
			return 0, false
		}
		value = -int64(v)
	} else {
		if v > math.MinInt64 {
			return 0, false
		}
		value = int64(v)
	}
	flag = true
	return
}

/* Convert a string into a long. Returns 1 if the string could be parsed into a
 * (non-overflowing) long, 0 otherwise. The value will be set to the parsed
 * value when appropriate. */
func string2l(s string) (long int32, flag bool) {
	v, flag := string2ll(s)
	if !flag {
		return 0, false
	}
	if v < math.MinInt32 || v > math.MaxInt32 {
		return 0, false
	}
	return int32(v), true
}

//int string2l(const char *s, size_t slen, long *lval) {
//long long llval;
//
//if (!string2ll(s, slen, &llval))
//return 0;
//
//if (llval < LONG_MIN || llval > LONG_MAX)
//return 0;
//
//*lval = (long)llval;
//return 1;
//}
