package types

import "sort"

type YamlNode struct {
	// 父节点
	Parent *YamlNode
	// 只有parent为null时候，该值才可能有值
	ProjectRemark string
	// name
	Name string
	// value
	Value string
	// 子节点
	Children []YamlNode
	// 数组标示
	ArrayFlag bool
	// 存储的数组中的前一个节点的下标
	LastNodeIndex int
	// 只有数组标示为true，下面的value才有值
	ValueList []YamlNode
}

func (yamlNode *YamlNode) ResortValue() {
	if !yamlNode.ArrayFlag || len(yamlNode.ValueList) == 0 {
		return
	}

	sort.Slice(yamlNode.ValueList, func(i, j int) bool {
		a := yamlNode.ValueList[i]
		b := yamlNode.ValueList[j]

		if -1 == a.LastNodeIndex || -1 == b.LastNodeIndex {
			return false
		}
		return a.LastNodeIndex < b.LastNodeIndex
	})

	for _, node := range yamlNode.ValueList {
		node.ResortValue()
	}
}
