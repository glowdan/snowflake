package main

import (
	"fmt"
	"time"
)

var second, workerId, wid int64 = -1, 0, -1
var id = -1

const MAX = 2<<15 - 8

func init() {
	second = time.Now().Unix()
	id = 0
	wid = nextWorkerID()
}

func main() {
	now := time.Now().UnixNano()
	for i := 0; i < 100000000; i++ {
		get()
	}
	fmt.Println((time.Now().UnixNano() - now) / 1000000)
	fmt.Println(workerId, wid, id)
}

func get() int64 {
	if id%1024 == 0 { //间断性获取时间，减少系统API调用
		now := time.Now().Unix()
		if now != second { //重新从0开始
			id = 0
		}
		if now < second { //时间回调，变更workerID
			wid = nextWorkerID()
			id = 0
		}
		second = now
	} else if id > MAX { //超限，换workID
		wid = nextWorkerID()
		id = 0
	}

	uniqueId := second<<31 + int64(id)<<15 + wid

	id += 1
	return uniqueId
}

func nextWorkerID() int64 {
	workerId += 1
	return workerId % (2 << 14)
}

