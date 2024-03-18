package jsonx

import "hash/fnv"

type hashCode uint32

func hashString(s string) hashCode {
	h := fnv.New32a()
	h.Write([]byte(s))
	return hashCode(h.Sum32())
}
