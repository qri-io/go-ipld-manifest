package manifest

import (
	"context"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-format"
)

// Manifest is a DAG of only block names and links (no content)
// node identifiers are stored in a slice "nodes", all other slices reference
// cids by index positions
type Manifest struct {
	Nodes []string `json:"nodes"`
	Links [][2]int `json:"links"`
	Sizes []uint64 `json:"sizes"`
}

// Node is a subset of the ipld format.Node interface
type Node interface {
	// pulled from blocks.Block format
	Cid() *cid.Cid
	// Links is a helper function that returns all links within this object
	Links() []*format.Link
	// Size returns the size in bytes of the serialized object
	Size() (uint64, error)
}

// NewManifest generates a manifest from an ipld node
func NewManifest(ctx context.Context, ng format.NodeGetter, node Node) (*Manifest, error) {
	ms := &mstate{
		ctx:  ctx,
		ng:   ng,
		cids: map[string]int{},
		m:    &Manifest{},
	}

	if _, err := ms.addNode(node); err != nil {
		return nil, err
	}
	return ms.m, nil
}

// mstate is a state machine for generating a manifest
type mstate struct {
	ctx  context.Context
	ng   format.NodeGetter
	idx  int
	cids map[string]int // lookup table of already-added cids
	m    *Manifest
}

// addNode places a node in the manifest & state machine, recursively adding linked nodes
// addNode returns early if this node is already added to the manifest
func (ms *mstate) addNode(node Node) (int, error) {
	id := node.Cid().String()

	if idx, ok := ms.cids[id]; ok {
		return idx, nil
	}

	// add the node
	idx := ms.idx
	ms.idx++

	ms.cids[id] = idx
	ms.m.Nodes = append(ms.m.Nodes, id)

	// ignore size errors b/c uint64 has no way to represent
	// errored size state as an int (-1), hopefully implementations default to 0
	// when erroring :/
	size, _ := node.Size()

	ms.m.Sizes = append(ms.m.Sizes, size)

	for _, link := range node.Links() {
		linkNode, err := link.GetNode(ms.ctx, ms.ng)
		if err != nil {
			return -1, err
		}

		nodeIdx, err := ms.addNode(linkNode)
		if err != nil {
			return -1, err
		}

		ms.m.Links = append(ms.m.Links, [2]int{idx, nodeIdx})
	}

	return idx, nil
}
