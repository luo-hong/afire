package gid

import "sync"

var syncMap sync.Map

func GetGidMap() *sync.Map {
	return &syncMap
}
