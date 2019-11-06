package redis


type dictEntry struct {
	key interface{}
	val interface{}
	next *dictEntry
}
type dict struct {
	ht map[*robj]*robj
}
