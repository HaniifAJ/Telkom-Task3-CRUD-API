package main

type HashTable struct {
	data [][][]string
}

func NewHashTable(size int) *HashTable {
	return &HashTable{
		data: make([][][]string, size),
	}
}

func (h *HashTable) _hash(key string) int {
	hash := 0
	for i := 0; i < len(key); i++ {
		hash = (hash + int(key[i])*i) % len(h.data)
	}
	return hash
}

func (h *HashTable) set(key string, value string) [][][]string {
	address := h._hash(key)
	if h.data[address] == nil {
		h.data[address] = make([][]string, 0)
	}
	h.data[address] = append(h.data[address], []string{key, value})
	return h.data
}

func (h *HashTable) get(key string) string {
	address := h._hash(key)
	currentBucket := h.data[address]
	if currentBucket != nil {
		for i := 0; i < len(currentBucket); i++ {
			if currentBucket[i][0] == key {
				return currentBucket[i][1]
			}
		}
	}
	return ""
}
