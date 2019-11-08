package redis

import "fmt"

/* -----------------------------------------------------------------------------
 * Higher level functions to queue data on the client output buffer.
 * The following functions are the ones that commands implementations will call.
 * -------------------------------------------------------------------------- */
func addReply(c *client, obj *robj) {
	// todo
	addReplyString(c, obj.getString())
}

func addReplyDouble(c *client, d float64) {
	dbuf := fmt.Sprintf("%.17g", d)
	sbuf := fmt.Sprintf("$%d\r\n%s\r\n", len(dbuf), dbuf)
	addReplyString(c, sbuf)
}

var replyBuf string

func addReplyString(c *client, s string) {
	// todo
	replyBuf = s
	return
}