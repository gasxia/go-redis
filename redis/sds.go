package redis

import "strconv"

type sds struct {
	len  int
	free int
	buf  []byte
}

func newsds(init string) *sds {
	buf := []byte(init)
	sh := &sds{
		len:  len(buf),
		free: 0,
		buf:  buf,
	}
	return sh
}

func (s *sds) sdslen() int {
	return s.len
}

func sdsfromlonglong(value int64) *sds {
	s := strconv.FormatInt(value, 10)
	return newsds(s)
}
