package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	epoch             = int64(1704067200000) // 2024-01-01 00:00:00 UTC
	nodeBits          = uint(10)
	stepBits          = uint(12)
	nodeMax           = int64(-1 ^ (-1 << nodeBits))
	stepMask          = int64(-1 ^ (-1 << stepBits))
	timeShift         = nodeBits + stepBits
	nodeShift         = stepBits
	maxClockBackwards = 5 * time.Millisecond
)

var (
	ErrInvalidNode    = errors.New("node number out of range")
	ErrClockBackwards = errors.New("system clock moved backwards beyond tolerance")
)

// Node Snowflake 发号节点
type Node struct {
	mu   sync.Mutex
	time int64
	node int64
	step int64
}

// NewNode 创建发号节点实例，node 范围 0 ~ 1023
func NewNode(node int64) (*Node, error) {
	if node < 0 || node > nodeMax {
		return nil, ErrInvalidNode
	}
	return &Node{
		time: 0,
		node: node,
		step: 0,
	}, nil
}

// Generate 生成 64 位全局唯一单调递增 ID
func (n *Node) Generate() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Now().UnixMilli()

	if now < n.time {
		diff := time.Duration(n.time-now) * time.Millisecond
		if diff <= maxClockBackwards {
			time.Sleep(diff)
			now = time.Now().UnixMilli()
		} else {
			for now < n.time {
				time.Sleep(time.Millisecond)
				now = time.Now().UnixMilli()
			}
		}
	}

	if n.time == now {
		n.step = (n.step + 1) & stepMask
		if n.step == 0 {
			for now <= n.time {
				time.Sleep(100 * time.Microsecond)
				now = time.Now().UnixMilli()
			}
		}
	} else {
		n.step = 0
	}

	n.time = now

	return ((now - epoch) << timeShift) | (n.node << nodeShift) | n.step
}

// DefaultNode 进程内默认单例发号节点
var defaultNode, _ = NewNode(1)

// GenerateID 获取全局唯一分布式 ID
func GenerateID() int64 {
	return defaultNode.Generate()
}
