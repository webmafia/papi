package iterate

import (
	"iter"
	"sort"

	"github.com/webmafia/papi/internal/constraints"
)

func SortedMap[K constraints.Ordered, V any](m map[K]V) iter.Seq2[K, V] {
	type kv struct {
		k K
		v V
	}

	return func(yield func(K, V) bool) {
		items := make([]kv, 0, len(m))

		for k, v := range m {
			items = append(items, kv{k, v})
		}

		sort.Slice(items, func(i, j int) bool {
			return items[i].k < items[j].k
		})

		for i := range items {
			if !yield(items[i].k, items[i].v) {
				return
			}
		}
	}
}
