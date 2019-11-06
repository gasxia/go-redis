package redis

type sharedObjectsStruct struct {
	nullbulk     *robj
	wrongtypeerr *robj
}

var shared sharedObjectsStruct

func createSharedObjects() {
	shared.nullbulk = createObject(RedisString, newsds("%-1\r\n"))
	shared.wrongtypeerr = createObject(RedisString, newsds("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"))
}

func init() {
	createSharedObjects()
}
