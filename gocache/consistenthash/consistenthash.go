package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 定义函数类型 Hash
type Hash func(data []byte) uint32

// Map constains all hashed keys
type Map struct {
	hash     Hash           // hash 函数
	replicas int            // 虚拟节点数量
	keys     []int          // hash , 有序
	hashMap  map[int]string // 虚拟节点和真实节点的映射表
}

// 实例化 Map, 允许自定义 hash 函数
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		// hash 函数默认为 crc32.ChecksumIEEE 算法
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash.
// 添加真实节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 对每一个真实节点 key，对应创建 m.replicas 个虚拟节点
		for i := 0; i < m.replicas; i++ {
			// 虚拟节点的名称是：strconv.Itoa(i) + key，即通过添加编号的方式区分不同虚拟节点
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash) // 把虚拟节点添加到环上
			m.hashMap[hash] = key         // 添加虚拟节点和真实节点的映射
		}
	}
	sort.Ints(m.keys) // 排序
}

// 根据 key 获取哈希环中最接近的节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	// Binary search for appropriate replica.
	// 在哈希环上二分查找虚拟节点对应的索引
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 返回对应的真实节点
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
