package crypto

import (
	"errors"

	. "github.com/elastos/Elastos.ELA.Utility/common"
)

type MerkleTree struct {
	Depth uint
	Root  *MerkleTreeNode
}

type MerkleTreeNode struct {
	Hash  Uint256
	Left  *MerkleTreeNode
	Right *MerkleTreeNode
}

func (t *MerkleTreeNode) IsLeaf() bool {
	return t.Left == nil && t.Right == nil
}

//use []Uint256 to create a new MerkleTree
func NewMerkleTree(hashes []Uint256) (*MerkleTree, error) {
	if len(hashes) == 0 {
		return nil, errors.New("NewMerkleTree input no item error.")
	}

	var height uint = 1
	nodes := generateLeaves(hashes)
	for len(nodes) > 1 {
		nodes = levelUp(nodes)
		height += 1
	}
	mt := &MerkleTree{
		Root:  nodes[0],
		Depth: height,
	}
	return mt, nil

}

//Generate the leaves nodes
func generateLeaves(hashes []Uint256) []*MerkleTreeNode {
	var leaves []*MerkleTreeNode
	for _, d := range hashes {
		node := &MerkleTreeNode{
			Hash: d,
		}
		leaves = append(leaves, node)
	}
	return leaves
}

//calc the next level's hash use double sha256
func levelUp(nodes []*MerkleTreeNode) []*MerkleTreeNode {
	var nextLevel []*MerkleTreeNode
	for i := 0; i < len(nodes)/2; i++ {
		hash, _ := MakeMerkleParent(&nodes[i*2].Hash, &nodes[i*2+1].Hash)
		node := &MerkleTreeNode{
			Hash:  *hash,
			Left:  nodes[i*2],
			Right: nodes[i*2+1],
		}
		nextLevel = append(nextLevel, node)
	}
	if len(nodes)%2 == 1 {
		hash, _ := MakeMerkleParent(&nodes[len(nodes)-1].Hash, nil)
		node := &MerkleTreeNode{
			Hash:  *hash,
			Left:  nodes[len(nodes)-1],
			Right: nodes[len(nodes)-1],
		}
		nextLevel = append(nextLevel, node)
	}
	return nextLevel
}

//input a []uint256, create a MerkleTree & calc the root hash
func ComputeRoot(hashes []Uint256) (Uint256, error) {
	if len(hashes) == 0 {
		return Uint256{}, errors.New("NewMerkleTree input no item error.")
	}
	if len(hashes) == 1 {
		return hashes[0], nil
	}
	tree, _ := NewMerkleTree(hashes)
	return tree.Root.Hash, nil
}