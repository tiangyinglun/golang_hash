package main

import (
	"hash/crc32"
	"sync"
	"sort"
	"strconv"
	"fmt"
	"math/rand"
)

const default_virtual_node_num = 100

type HashRing []uint32
type Node struct {
	Ip   string
	Port string
}

type Consistent struct {
	Nodes   map[uint32]*Node             //所有虚拟节点 map[ssss]={}
	numReps int                          //数量
	ipMap   map[string]map[uint32]string //
	ring    HashRing                     //hash表
	sync.RWMutex
}

type hashInterface interface {
	AddNode(arKey string) bool
	GetNode(arKey string) string
	Remove(arKey string)
}

func main() {
	c := NewConsistent()
	var ha hashInterface
	ha = c
	for i := 0; i < 5; i++ {

		str := "127.0.0." + strconv.Itoa(i)
		ha.AddNode(str)
	}
	mc := make(map[string]int)
	for i := 0; i < 1000; i++ {
		it:=rand.Int()
		str := "127.0.0." + strconv.Itoa(it)
		m := ha.GetNode(str)
		mc[m] = mc[m] + 1
	}
	fmt.Println(mc)

}

/**
 新节点
 */
func NewNode(ip, port string) *Node {
	return &Node{
		Ip:   ip,
		Port: port,
	}
}

func (c HashRing) Len() int {
	return len(c)
}

func (c HashRing) Less(i, j int) bool {
	return c[i] < c[j]
}

func (c HashRing) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

/**
初始化
 */
func NewConsistent() *Consistent {
	return &Consistent{
		Nodes:   make(map[uint32]*Node),
		numReps: default_virtual_node_num,
		ipMap:   make(map[string]map[uint32]string),
		ring:    HashRing{},
	}
}

/**
 添加节点
 */
func (c *Consistent) AddNode(arKey string) bool {
	c.Lock()
	defer c.Unlock()

	for i := 0; i < c.numReps; i++ {
		mp := make(map[uint32]string)
		if i == 0 {
			hashCrc := hashKey(arKey)  //crc
			node := NewNode(arKey, "") //crc 节点
			c.Nodes[hashCrc] = node

			mp[hashCrc] = arKey
		} else {
			hashC := hashKey(arKey + strconv.Itoa(i))
			cnode := NewNode(arKey, "") //crc 节点
			c.Nodes[hashC] = cnode
			mp[hashC] = arKey
		}

		c.ipMap[arKey] = mp
	}
	c.SortHashRing()
	return true
}

func (c *Consistent) GetNode(arKey string) string {
	c.Lock()
	defer c.Unlock()
	hashCrc := hashKey(arKey)

	if len(c.ring) == 0 {
		return ""
	}
	for _, v := range c.ring {
		if hashCrc <= v {
			return c.Nodes[v].Ip
		}
	}
	return c.Nodes[c.ring[len(c.ring)-1]].Ip

}

//删除节点
func (c *Consistent) Remove(arKey string) {
	c.Lock()
	defer c.Unlock()
	hashCrc := hashKey(arKey)
	_, ok := c.Nodes[hashCrc]
	if !ok {
		return
	}
	for k, _ := range c.ipMap[arKey] {
		delete(c.Nodes, k)
	}
	delete(c.ipMap, arKey)
	c.SortHashRing()

}

func (c *Consistent) SortHashRing() {
	c.ring = HashRing{}
	for k := range c.Nodes {
		c.ring = append(c.ring, k)
	}
	sort.Sort(c.ring)
}

/**
  生成散列值
 */
func hashKey(arKey string) (u uint32) {
	arKeyByte := []byte(arKey)
	u = crc32.ChecksumIEEE(arKeyByte)
	return
}
