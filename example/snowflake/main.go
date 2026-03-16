package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	workerBits  uint8 = 10                      //机器码位数
	numberBits  uint8 = 12                      //序列号位数
	workerMax   int64 = -1 ^ (-1 << workerBits) //机器码最大值（即1023）
	numberMax   int64 = -1 ^ (-1 << numberBits) //序列号最大值（即4095）
	timeShift         = workerBits + numberBits //时间戳偏移量
	workerShift       = numberBits              //机器码偏移量
	epoch       int64 = 828838800000            //起始常量时间戳（毫秒)
)

type Worker struct {
	mu        sync.Mutex
	timeStamp int64
	workerId  int64
	number    int64
}

// NewWorker 生成一个新的工作节点,
// 参数为当前节点的ID,最大为1023,
// 如果超过了最大值,则会返回错误
func NewWorker(workerId int64) (*Worker, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("WorkerId超过了限制！")
	}
	return &Worker{
		timeStamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}

func (w *Worker) NextId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	//当前时间的毫秒时间戳
	now := time.Now().UnixNano() / 1e6
	//如果时间戳与当前时间相同，则增加序列号
	if w.timeStamp == now {
		w.number++
		//如果序列号超过了最大值，则更新时间戳
		if w.number > numberMax {
			for now <= w.timeStamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else { //如果时间戳与当前时间不同，则直接更新时间戳
		w.number = 0
		w.timeStamp = now
	}
	//ID由时间戳、机器编码、序列号组成
	ID := (now-epoch)<<timeShift | (w.workerId << workerShift) | (w.number)
	return ID
}

func main() {
	fid, _ := NewWorker(0)
	id := fid.NextId() // 3960567221469577216
	fmt.Printf("%064b\n", id)  // 补齐64位二进制
	fmt.Printf("%d\n", id)     // 普通十进制打印（省略前导0）
}
