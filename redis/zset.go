package redis

/*
 * 有序集合
 */
type zset struct {
	// 字典，键为成员，值为分值
	// 用于支持 O(1) 复杂度的按成员取分值操作
	dict *dict

	zsl *skiplist
}

func zaddGenericCommand(c *client, incr int) {

}

func zaddCommand(c *client) {
	zaddGenericCommand(c, 0)
}

func zincrbyCommand(c *client) {
	zaddGenericCommand(c, 1)
}

func zscoreCommand(c *client) {
	key := c.argv[1]
	//var score float64

	zobj := lookipKeyReadOrReply(c, key, shared.nullbulk)
	if zobj == nil || zobj.checkType(c, RedisZset) {
		return
	}

	if zobj.encoding == RedisEncodingZiplist {
		// todo
	} else if zobj.encoding == RedisEncodingSkiplist {
		zs := zobj.ptr.(*zset)

		c.argv[2] = tryObjectEncoding(c.argv[2])
		de := zs.dict.Find(c.argv[2])
		if de != nil {
			if score, ok := de.GetVal().(float64); ok {
				addReplyDouble(c, score)
			}
		} else {
			addReply(c, shared.nullbulk)
		}
	} else {
		panic("Unknown sorted set encoding")
	}

}
