package gocache

import pb "go-cache/gocachepb"

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.
type PeerPicker interface {
	// 用于根据传入的 key 选择相应节点 PeerGetter
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer.
type PeerGetter interface {
	// 从对应 group 查找缓存值
	Get(in *pb.Request, out *pb.Response) error
}
