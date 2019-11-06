package redis

import "strings"

/* A redis object, that is a type able to hold a string / list / set */

/* The actual Redis Object */
type redisObjectType uint8
type redisObjectEncoding uint8

const (
	RedisString redisObjectType = iota
	RedisList
	RedisSet
	RedisZset
	RedisHash
)

const (
	RedisEncodingRaw        redisObjectEncoding = iota /* Raw representation */
	RedisEncodingInt                                   /* Encoded as integer */
	RedisEncodingHt                                    /* Encoded as hash table */
	RedisEncodingZipmap                                /* Encoded as zipmap */
	RedisEncodingLinkedlist                            /* Encoded as regular linked list */
	RedisEncodingZiplist                               /* Encoded as ziplist */
	RedisEncodingIntset                                /* Encoded as intset */
	RedisEncodingSkiplist                              /* Encoded as skiplist */
	RedisEncodingEmbstr                                /* Embedded sds string encoding */
)

type redisObject struct {
	_type    redisObjectType
	encoding redisObjectEncoding
	lru      uint32 /* lru time (relative to server.lruclock) */
	refcount int
	ptr      interface{}
}

type robj redisObject

func compareStringObjects(obj *robj, other *robj) int {
	return strings.Compare(obj.getString(), other.getString())
	//return obj.getString() == other.getString()
}

func (obj *robj) getString() string {
	return obj.ptr.(string)
}

/* Equal string objects return 1 if the two objects are the same from the
 * point of view of a string comparison, otherwise 0 is returned.
 *
 * 如果两个对象的值在字符串的形式上相等，那么返回 1 ， 否则返回 0 。
 *
 * Note that this function is faster then checking for (compareStringObject(a,b) == 0)
 * because it can perform some more optimization.
 *
 * 这个函数做了相应的优化，所以比 (compareStringObject(a, b) == 0) 更快一些。
 */
func equalStringObjects(a *robj, b *robj) bool {
	// 对象的编码为 INT ，直接对比值
	// 这里避免了将整数值转换为字符串，所以效率更高
	/* If both strings are integer encoded just check if the stored
	 * long is the same. */
	if a.encoding == RedisEncodingInt && b.encoding == RedisEncodingInt {
		return a.ptr == b.ptr
	} else {
		return compareStringObjects(a, b) == 0
	}
}

func createObject(rtype redisObjectType, ptr interface{}) *robj {
	return &robj{
		_type:    rtype,
		encoding: RedisEncodingRaw,
		lru:      1,
		refcount: 1,
		ptr:      ptr,
	}
}

/* Create a string object with encoding REDIS_ENCODING_RAW, that is a plain * string object where o->ptr points to a proper sds string. */
// 创建一个 REDIS_ENCODING_RAW 编码的字符对象
// 对象的指针指向一个 sds 结构
func createRawStringObject(ptr string) *robj {
	return createObject(RedisString, newsds(ptr))
}

/* Create a string object with encoding REDIS_ENCODING_EMBSTR, that is
 * an object where the sds string is actually an unmodifiable string
 * allocated in the same chunk as the object itself. */
// 创建一个 REDIS_ENCODING_EMBSTR 编码的字符对象
// 这个字符串对象中的 sds 会和字符串对象的 redisObject 结构一起分配
// 因此这个字符也是不可修改的

// jrs: golang中string类型也是
func createEmbeddedStringObject(ptr string) *robj {
	return &robj{
		_type:    RedisString,
		encoding: RedisEncodingEmbstr,
		lru:      1,
		refcount: 1,
		ptr:      ptr,
	}
	//return createRawStringObject(ptr)
}

/* Create a string object with EMBSTR encoding if it is smaller than
 * REIDS_ENCODING_EMBSTR_SIZE_LIMIT, otherwise the RAW encoding is
 * used.
 *
 * The current limit of 39 is chosen so that the biggest string object
 * we allocate as EMBSTR will still fit into the 64 byte arena of jemalloc. */
const RedisEncodingEmbstrSizeLimit = 39

func createStringObject(ptr string) *robj {
	if len(ptr) <= RedisEncodingEmbstrSizeLimit {
		return createEmbeddedStringObject(ptr)
	} else {
		return createRawStringObject(ptr)
	}
}

/*
 * 检查对象 o 的类型是否和 type 相同：
 *
 *  - 相同返回 0
 *
 *  - 不相同返回 1 ，并向客户端回复一个错误
 */
func (obj *robj)checkType(c *client, _type redisObjectType) bool {
	if obj._type != _type {
		addReply(c,shared.wrongtypeerr)
		return false
	}
	return true
}

func (obj *robj)sdsEncodedObject() bool {
	return obj.encoding == RedisEncodingRaw || obj.encoding == RedisEncodingEmbstr
}

func tryObjectEncoding(o *robj) *robj {
	var value int64
	s := o.ptr.(*sds)

	// todo assert
	// redisAssertWithInfo...

	if !o.sdsEncodedObject() {
		return o
	}

	len := s.sdslen()
	if len <= 21 && string2l(s, len) {

	}


}
