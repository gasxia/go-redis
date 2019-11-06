package redis

type client struct {
	db   *Db
	argv []*robj
}
