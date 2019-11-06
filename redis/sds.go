package redis

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