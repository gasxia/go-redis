package redis

type Db struct {
	id int
}

/*
 * 为执行读取操作而取出键 key 在数据库 db 中的值。
 *
 * 并根据是否成功找到值，更新服务器的命中/不命中信息。
 *
 * 找到时返回值对象，没找到返回 NULL 。
 */
func lookupKeyRead(db *Db, key *robj) *robj {
	// todo
	return nil
	//robj *val;
	//
	//// 检查 key 释放已经过期
	//expireIfNeeded(db,key);
	//
	//// 从数据库中取出键的值
	//val = lookupKey(db,key);
	//
	//// 更新命中/不命中信息
	//if (val == NULL)
	//server.stat_keyspace_misses++;
	//else
	//server.stat_keyspace_hits++;
	//
	//// 返回值
	//return val;
}

/*
 * 为执行读取操作而从数据库中查找返回 key 的值。
 *
 * 如果 key 存在，那么返回 key 的值对象。
 *
 * 如果 key 不存在，那么向客户端发送 reply 参数中的信息，并返回 NULL 。
 */
func lookipKeyReadOrReply(c *client, key *robj, reply *robj) *robj {
	o := lookupKeyRead(c.db, key)
	if o == nil {
		addReply(c, reply)
	}
	return o
}
