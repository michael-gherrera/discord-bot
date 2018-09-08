package util

type MaxStockHeap struct {
	MinStockHeap
}

func (h MaxStockHeap) Less(i, j int) bool {
	return h.MinStockHeap[i].LDCost > h.MinStockHeap[j].LDCost
}
