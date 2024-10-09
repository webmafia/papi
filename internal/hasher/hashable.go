package hasher

type Hashable interface {
	Hash() uint64
}
