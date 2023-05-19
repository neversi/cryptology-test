package depth_test

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"stock/updater/depth"
	"stock/updater/depth/bs"
	"stock/updater/depth/simple"
	"stock/updater/models"

	"github.com/google/go-cmp/cmp"
)

var stock depth.Depth

func TestMain(m *testing.M) {
	fmt.Println("here")
	stocks := map[string]depth.Depth{
		"Simple Implementation":        simple.NewDepth(),
		"Binary Search Implementation": bs.NewDepth(depth.Options{}),
	}
	var code int
	rand.Seed(time.Now().UnixNano())
	for name, impl := range stocks {
		fmt.Println(strings.Repeat("–", 30))
		fmt.Println(name)
		fmt.Println(strings.Repeat("–", 30))
		stock = impl
		code = m.Run()
	}

	os.Exit(code)
}

func TestStoreAndLoad(t *testing.T) {
	testCases := map[string]struct {
		opts     depth.Options
		in       []models.Item
		expected []models.Item
	}{
		"sorted items for bids": {
			opts:     depth.Options{true},
			in:       []models.Item{{1.5, 1}, {1.3, 3}, {1.1, 3}},
			expected: []models.Item{{1.5, 1}, {1.3, 3}, {1.1, 3}},
		},
		"unsorted items for bids": {
			opts:     depth.Options{true},
			in:       []models.Item{{5.4, 3}, {3.9, 1}, {7.6, 5}},
			expected: []models.Item{{7.6, 5}, {5.4, 3}, {3.9, 1}},
		},
		"sorted items for asks": {
			opts:     depth.Options{false},
			in:       []models.Item{{1.1, 3}, {1.3, 3}, {1.5, 1}},
			expected: []models.Item{{1.1, 3}, {1.3, 3}, {1.5, 1}},
		},
		"unsorted items for asks": {
			opts:     depth.Options{false},
			in:       []models.Item{{5.4, 3}, {3.9, 1}, {7.6, 5}},
			expected: []models.Item{{3.9, 1}, {5.4, 3}, {7.6, 5}},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			stock.SetOptions(tc.opts)
			stock.Store(tc.in)
			if !cmp.Equal(stock.Load(), tc.expected) {
				t.Fatalf("invalid order of items:\n%v", cmp.Diff(stock.Load(), tc.expected))
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	testCases := map[string]struct {
		opts     depth.Options
		initial  []models.Item
		expected []models.Item
		in       models.Item
	}{
		"insert new item in bids": {
			opts:     depth.Options{true},
			initial:  []models.Item{{1.5, 1}, {1.3, 3}, {1.1, 3}},
			in:       models.Item{1.6, 2},
			expected: []models.Item{{1.6, 2}, {1.5, 1}, {1.3, 3}, {1.1, 3}},
		},
		"delete item from bids": {
			opts:     depth.Options{true},
			initial:  []models.Item{{1.5, 1}, {1.3, 3}, {1.1, 3}},
			in:       models.Item{1.1, 0},
			expected: []models.Item{{1.5, 1}, {1.3, 3}},
		},
		"modify item from bids": {
			opts:     depth.Options{true},
			initial:  []models.Item{{1.5, 1}, {1.3, 3}, {1.1, 3}},
			in:       models.Item{1.3, 2},
			expected: []models.Item{{1.5, 1}, {1.3, 2}, {1.1, 3}},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			stock.SetOptions(tc.opts)
			stock.Store(tc.initial)
			stock.Update(tc.in)

			if !cmp.Equal(tc.expected, stock.Load()) {
				t.Fatalf("invalid insertion of order:\n%v", cmp.Diff(stock.Load(), tc.expected))
			}
		})
	}
}

// BenchmarkUpdate benches update method of several implementations
// binary search implementation in this case not so efficient as amortized
// insertion of element is done in O(N) meanwhile, appending in simple implementation
// we append to the end (which is O(1) amortized) and eventually sort it (with the best
// case as our array is already sorted)
// binary search could be realized through BST which will give more efficiency in terms
// of modification but less in terms of representation
func BenchmarkUpdate(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	generateItem := func() models.Item {
		return models.Item{
			Price: float64(int(rand.Float64()*100)) / 100,
			Size:  1,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var item models.Item
		item = generateItem() // const time negligble in benchmarks (as stop/start timer takes more time to execute)
		stock.Update(item)
	}
}
