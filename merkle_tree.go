package main

import "crypto/sha256"

type MerkleTreeNode struct {
	Data []byte
	Left *MerkleTreeNode
	Right *MerkleTreeNode
}

type MerkleTree struct {
	Node *MerkleTreeNode
}

func NewMerkleNode(left,right *MerkleTreeNode, data []byte)  *MerkleTreeNode{
	node := MerkleTreeNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	}else {
		prevHash := append(left.Data,right.Data...)
		hash := sha256.Sum256(prevHash)
		node.Data = hash[:]
	}
	node.Left = left
	node.Right = right
	return &node
}

func NewMerkleTree(data [][]byte) *MerkleTree{
	var nodes []MerkleTreeNode

	if len(data)%2 != 0{ //如果节点数量不是偶数的话，那么增加最后一个以凑成偶数
		data = append(data,data[len(data)-1])
	}
	for _, nodeData := range data{
		node := MerkleTreeNode{nodeData,nil,nil}
		nodes = append(nodes,node)
	}

	for i := 0; i < len(data)/2; i++ {
		var newLevel []MerkleTreeNode

		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}

		nodes = newLevel
	}

	mTree := MerkleTree{&nodes[0]}

	return &mTree
}