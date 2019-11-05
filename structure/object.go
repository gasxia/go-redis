package structure

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
	rtype    redisObjectType
	encoding redisObjectEncoding
	lru      uint32 /* lru time (relative to server.lruclock) */
	refcount int
	ptr      interface{}
}

type robj redisObject

func (obj *robj) compareStringObjects(other *robj) bool {
	return obj.getString() == other.getString()
}

func (obj *robj) getString() string {
	return obj.ptr.(string)
}

func newObject(rtype redisObjectType, ptr interface{}) *robj {
	return &robj{
		rtype:    rtype,
		encoding: RedisEncodingRaw,
		lru:      1,
		refcount: 1,
		ptr:      ptr,
	}
}

/* Create a string object with encoding REDIS_ENCODING_RAW, that is a plain * string object where o->ptr points to a proper sds string. */
// 创建一个 REDIS_ENCODING_RAW 编码的字符对象
// 对象的指针指向一个 sds 结构
func newRawStringObject(ptr string) *robj {
	return newObject(RedisString, newsds(ptr))
}

/* Create a string object with encoding REDIS_ENCODING_EMBSTR, that is
 * an object where the sds string is actually an unmodifiable string
 * allocated in the same chunk as the object itself. */
// 创建一个 REDIS_ENCODING_EMBSTR 编码的字符对象
// 这个字符串对象中的 sds 会和字符串对象的 redisObject 结构一起分配
// 因此这个字符也是不可修改的

// jrs: golang中string类型也是
func newEmbeddedStringObject(ptr string) *robj {
	return &robj{
		rtype:    RedisString,
		encoding: RedisEncodingEmbstr,
		lru:      1,
		refcount: 1,
		ptr:      ptr,
	}
	//return newRawStringObject(ptr)
}

/* Create a string object with EMBSTR encoding if it is smaller than
 * REIDS_ENCODING_EMBSTR_SIZE_LIMIT, otherwise the RAW encoding is
 * used.
 *
 * The current limit of 39 is chosen so that the biggest string object
 * we allocate as EMBSTR will still fit into the 64 byte arena of jemalloc. */
const RedisEncodingEmbstrSizeLimit = 39

func newStringObject(ptr string) *robj {
	if len(ptr) <= RedisEncodingEmbstrSizeLimit {
		return newEmbeddedStringObject(ptr)
	} else {
		return newRawStringObject(ptr)
	}
}

//#define REDIS_LRU_BITS 24
//#define REDIS_LRU_CLOCK_MAX ((1<<REDIS_LRU_BITS)-1) /* Max value of obj->lru */
//#define REDIS_LRU_CLOCK_RESOLUTION 1000 /* LRU clock resolution in ms */
//typedef struct redisObject {
//
//// 类型
//unsigned type:4;
//
//// 编码
//unsigned encoding:4;
//
//// 对象最后一次被访问的时间
//unsigned lru:REDIS_LRU_BITS;
//
//// 引用计数
//int refcount;
//
//// 指向实际值的指针
//void *ptr;
//
//} robj;
