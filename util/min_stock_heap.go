package util

type MinStockHeap []*struct {
	// Levenshtein Distance Cost
	LDCost int
	Symbol string
}

func (h MinStockHeap) Len() int           { return len(h) }
func (h MinStockHeap) Less(i, j int) bool { return h[i].LDCost < h[j].LDCost }
func (h MinStockHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinStockHeap) Push(x interface{}) {
	*h = append(*h, x.(*struct {
		LDCost int
		Symbol string
	}))
}

func (h *MinStockHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
