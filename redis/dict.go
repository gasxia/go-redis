package redis


type dictEntry struct {
	key interface{}
	val interface{}
	next *dictEntry
}

func (de *dictEntry) GetVal() interface{} {
	return de.val
}

type dict struct {
	ht map[interface{}]interface{}
}

func (dt *dict) Find(key *robj) *dictEntry {
	if val, ok := dt.ht[key]; ok {
		return &dictEntry{
			key: key,
			val: val,
		}
	}
	return nil
}
