package main

// 引入需要的包
import (
	"fmt"
	"hash/crc32"
	"math"
	"sort"
	"strconv"
)

// 定义虚拟节点数量
const VIRTUAL_NODE_NUM = 100

// 定义环的最大长度
const MAX_HASH_LEN = 1 << 31

// 定义数据节点结构体
type DataNode struct {
	Name string
}

// 定义虚拟节点结构体
type VirtualNode struct {
	Name      string
	DataNode  *DataNode
	HashValue uint32
}

// 定义一致性 Hash 结构体
type ConsistentHash struct {
	DataNodes       []*DataNode
	VirtualNodes    []*VirtualNode
	VirtualNodeNums map[string]int
}

// 实现一致性 Hash 结构体的排序接口
func (ch *ConsistentHash) Len() int {
	return len(ch.VirtualNodes)
}

func (ch *ConsistentHash) Less(i, j int) bool {
	return ch.VirtualNodes[i].HashValue < ch.VirtualNodes[j].HashValue
}

func (ch *ConsistentHash) Swap(i, j int) {
	ch.VirtualNodes[i], ch.VirtualNodes[j] = ch.VirtualNodes[j], ch.VirtualNodes[i]
}

// 创建一致性 Hash 结构体
func NewConsistentHash() *ConsistentHash {
	return &ConsistentHash{
		VirtualNodes:    make([]*VirtualNode, 0),
		VirtualNodeNums: make(map[string]int),
	}
}

// 添加数据节点
func (ch *ConsistentHash) AddDataNode(dataNode *DataNode) {
	// 将数据节点加入到数据节点列表中
	ch.DataNodes = append(ch.DataNodes, dataNode)
	// 为每个数据节点创建 VIRTUAL_NODE_NUM 个虚拟节点
	for i := 0; i < VIRTUAL_NODE_NUM; i++ {
		// 创建虚拟节点
		virtualNode := &VirtualNode{
			Name:      dataNode.Name + ":" + strconv.Itoa(i),
			DataNode:  dataNode,
			HashValue: ch.hash(virtualNode.Name),
		}
		// 将虚拟节点加入到虚拟节点列表中
		ch.VirtualNodes = append(ch.VirtualNodes, virtualNode)
	}
	// 对虚拟节点列表按照 Hash 值进行排序
	sort.Sort(ch)
	// 记录每个数据节点对应的虚拟节点数量
	ch.VirtualNodeNums[dataNode.Name] = VIRTUAL_NODE_NUM
}

// 删除数据节点
func (ch *ConsistentHash) RemoveDataNode(dataNode *DataNode) {
	// 从数据节点列表中删除数据节点
	for i := 0; i < len(ch.DataNodes); i++ {
		if ch.DataNodes[i].Name == dataNode.Name {
			ch.DataNodes = append(ch.DataNodes[:i], ch.DataNodes[i+1:]...)
			break
		}
	}
	// 从虚拟节点列表中删除该数据节点对应的所有虚拟节点
	for i := 0; i < len(ch.VirtualNodes); i++ {
		if ch.VirtualNodes[i].DataNode.Name == dataNode.Name {
			ch.VirtualNodes = append(ch.VirtualNodes[:i], ch.VirtualNodes[i+1:]...)
			i--
		}
	}
	// 更新记录每个数据节点对应的虚拟节点数量
	delete(ch.VirtualNodeNums, dataNode.Name)
}

// 获取指定数据项对应的数据节点
func (ch *ConsistentHash) GetDataNode(data string) *DataNode {
	// 如果没有数据节点，则返回 nil
	if len(ch.DataNodes) == 0 {
		return nil
	}
	// 计算数据项的 Hash 值
	hashValue := ch.hash(data)
	// 寻找第一个大于等于该 Hash 值的虚拟节点
	i := sort.Search(len(ch.VirtualNodes), func(i int) bool {
		return ch.VirtualNodes[i].HashValue >= hashValue
	})
	// 如果没有找到，则返回第一个虚拟节点对应的数据节点
	if i == len(ch.VirtualNodes) {
		return ch.VirtualNodes[0].DataNode
	}
	// 否则，返回找到的虚拟节点对应的数据节点
	return ch.VirtualNodes[i].DataNode
}

// 计算字符串的 Hash 值
func (ch *ConsistentHash) hash(data string) uint32 {
	return crc32.ChecksumIEEE([]byte(data)) % MAX_HASH_LEN
}

// 计算 KV 数据在服务器上分布数量的标准差
func (ch *ConsistentHash) StdDev() float64 {
	// 计算数据节点数量
	dataNodeNum := len(ch.DataNodes)
	// 计算每个数据节点对应的虚拟节点数量
	virtualNodeNums := make([]int, dataNodeNum)
	for i := 0; i < dataNodeNum; i++ {
		virtualNodeNums[i] = ch.VirtualNodeNums[ch.DataNodes[i].Name]
	}
	// 计算 KV 数据在服务器上分布数量的标准差
	return math.StdDev(virtualNodeNums)
}

// 主函数
func main() {
	// 创建一致性 Hash 结构体
	ch := NewConsistentHash()
	// 添加 10 个数据节点
	for i := 0; i < 10; i++ {
		ch.AddDataNode(&DataNode{
			Name: strconv.Itoa(i),
		})
	}
	// 测试 100 万个 KV 数据
	for i := 0; i < 1000000; i++ {
		// 获取数据项对应的数据节点
		dataNode := ch.GetDataNode(strconv.Itoa(i))
		// 在该数据节点对应的虚拟节点数量上加 1
		ch.VirtualNodeNums[dataNode.Name]++
	}
	// 计算 KV 数据在服务器上分布数量的标准差
	stdDev := ch.StdDev()
	fmt.Printf("KV数据在服务器上分布数量的标准差：%.2f\n", stdDev)
}

/*
定义虚拟节点数量：为了保证一致性 Hash 的性能，需要定义虚拟节点数量，并在添加数据节点时为每个数据节点创建相应数量的虚拟节点。
定义环的最大长度：在一致性 Hash 算法中，所有的虚拟节点会构成一个环，需要定义环的最大长度，并在计算 Hash 值时将其取模。
计算 Hash 值：可以使用 golang 中的哈希算法库来计算字符串的 Hash 值，例如可以使用 crc32.ChecksumIEEE() 来计算。
排序虚拟节点：在一致性 Hash 算法中，虚拟节点需要按照 Hash 值从小到大的顺序进行排序，可以使用 golang 中的 sort 库实现排序功能。
计算标准差：在一致性 Hash 算法中，需要计算 KV 数据在服务器上分布数量的标准差，可以使用 golang 中的 math 库实现计算标准差的功能。
*/

/*
通过按照上述步骤实现一致性 Hash 算法，可以在 golang 中实现一个可用的一致性 Hash 算法。此外，在实际应用中，还可以对一致性
 Hash 算法进行优化，提高算法的性能。例如，可以使用带有缓存的二分查找算法来提高查找虚拟节点的效率，或者使用环形链表来存储
 虚拟节点，提高添加和删除虚拟节点的效率。
*/
