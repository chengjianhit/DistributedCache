package consistent

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashSum func(data []byte) uint32

type ConsistentHash struct {
	hashSum    HashSum
	replicas   int
	hashCircle []int
	nodeMap    map[int]string //节点 hash值-节点名
}

func InitConsistentHashMap(replicas int, fnc HashSum) *ConsistentHash {
	consistentMap := &ConsistentHash{
		hashSum:  fnc,
		replicas: replicas,
		nodeMap:  make(map[int]string),
	}
	if fnc == nil {
		consistentMap.hashSum = crc32.ChecksumIEEE
	}
	return consistentMap
}

func (cnsHash *ConsistentHash) AddNode(nodeNames ...string) {
	for _, name := range nodeNames {
		for i := 0; i < cnsHash.replicas; i++ {
			hashSum := int(cnsHash.hashSum([]byte(strconv.Itoa(i) + name)))
			cnsHash.hashCircle = append(cnsHash.hashCircle, hashSum)
			cnsHash.nodeMap[hashSum] = name
		}
	}
	//对所有节点（真实节点、虚拟节点）的hash值进行排序，形成顺时针环
	sort.Ints(cnsHash.hashCircle)
}

// 获取距离该Key顺时针距离最近的节点
func (cnsHash *ConsistentHash) GetClosestNode(key string) string {

	keyHashSum := int(cnsHash.hashSum([]byte(key)))

	nodeIdx := sort.Search(len(cnsHash.hashCircle), func(i int) bool {
		return cnsHash.hashCircle[i] > keyHashSum
	})

	return cnsHash.nodeMap[cnsHash.hashCircle[nodeIdx%len(cnsHash.hashCircle)]]
}
