package lrustruct

type BaseStore struct {
	b []byte
}

func (baseStore BaseStore) Len() int {
	return len(baseStore.b)
}

func (baseStore BaseStore) CloneBaseStore() []byte {
	return cloneBytes(baseStore.b)
}

func (baseStore BaseStore) parseToString() string {
	return string(baseStore.b)
}

func cloneBytes(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(src, dst)
	return dst
}
