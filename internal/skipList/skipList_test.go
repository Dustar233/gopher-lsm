package skipList

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

// ======================== Constructor Tests ========================

func TestNew(t *testing.T) {
	sl := New(func(a, b any) int {
		return a.(int) - b.(int)
	})

	if sl == nil {
		t.Fatal("New returned nil")
	}
	if sl.level != 1 {
		t.Errorf("expected level=1, got %d", sl.level)
	}
	if sl.size != 0 {
		t.Errorf("expected size=0, got %d", sl.size)
	}
	if sl.maxLevel != defaultMaxLevel {
		t.Errorf("expected maxLevel=%d, got %d", defaultMaxLevel, sl.maxLevel)
	}
	if sl.P != defaultMaxP {
		t.Errorf("expected P=%d, got %d", defaultMaxP, sl.P)
	}
	if sl.headNode == nil {
		t.Fatal("headNode is nil")
	}
	if len(sl.headNode.nextNode) != defaultMaxLevel {
		t.Errorf("expected headNode.nextNode len=%d, got %d", defaultMaxLevel, len(sl.headNode.nextNode))
	}
}

func TestNewInt(t *testing.T) {
	sl := NewInt()
	if sl == nil {
		t.Fatal("NewInt returned nil")
	}
	if sl.size != 0 {
		t.Errorf("expected size=0, got %d", sl.size)
	}
}

// ======================== Set Tests ========================

func TestSet_InsertSingle(t *testing.T) {
	sl := NewInt()

	old, existed := sl.Set(1, "one")
	if existed {
		t.Error("Set new key should return existed=false")
	}
	if old != nil {
		t.Errorf("Set new key should return old=nil, got %v", old)
	}
	if sl.size != 1 {
		t.Errorf("expected size=1, got %d", sl.size)
	}
}

func TestSet_InsertMultiple(t *testing.T) {
	sl := NewInt()

	for i := 1; i <= 10; i++ {
		old, existed := sl.Set(i, fmt.Sprintf("val-%d", i))
		if existed {
			t.Errorf("Set new key %d should return existed=false", i)
		}
		if old != nil {
			t.Errorf("Set new key %d should return old=nil, got %v", i, old)
		}
	}

	if sl.size != 10 {
		t.Errorf("expected size=10, got %d", sl.size)
	}
}

func TestSet_UpdateExisting(t *testing.T) {
	sl := NewInt()

	sl.Set(1, "first")
	old, existed := sl.Set(1, "second")

	if !existed {
		t.Error("Set existing key should return existed=true")
	}
	if old != "first" {
		t.Errorf("expected old='first', got %v", old)
	}
	if sl.size != 1 {
		t.Errorf("expected size=1 after update, got %d", sl.size)
	}
}

func TestSet_UpdateExistingMultipleTimes(t *testing.T) {
	sl := NewInt()

	for i := 0; i < 5; i++ {
		old, existed := sl.Set(42, i)
		if i == 0 {
			if existed {
				t.Error("first insert should return existed=false")
			}
		} else {
			if !existed {
				t.Errorf("update #%d should return existed=true", i)
			}
			if old != i-1 {
				t.Errorf("update #%d: expected old=%d, got %v", i, i-1, old)
			}
		}
	}

	if sl.size != 1 {
		t.Errorf("expected size=1, got %d", sl.size)
	}
}

func TestSet_NilValue(t *testing.T) {
	sl := NewInt()

	old, existed := sl.Set(1, nil)
	if existed {
		t.Error("Set new key with nil value should return existed=false")
	}
	if old != nil {
		t.Errorf("expected old=nil, got %v", old)
	}
}

// ======================== Get Tests ========================

func TestGet_Exists(t *testing.T) {
	sl := NewInt()
	sl.Set(1, "one")
	sl.Set(2, "two")
	sl.Set(3, "three")

	val, ok := sl.Get(2)
	if !ok {
		t.Error("Get existing key should return ok=true")
	}
	if val != "two" {
		t.Errorf("expected 'two', got %v", val)
	}
}

func TestGet_NotExists(t *testing.T) {
	sl := NewInt()
	sl.Set(1, "one")

	val, ok := sl.Get(999)
	if ok {
		t.Error("Get non-existing key should return ok=false")
	}
	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}
}

func TestGet_EmptyList(t *testing.T) {
	sl := NewInt()

	val, ok := sl.Get(1)
	if ok {
		t.Error("Get on empty list should return ok=false")
	}
	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}
}

func TestGet_AllInsertedKeys(t *testing.T) {
	sl := NewInt()
	expected := map[int]string{
		5:  "five",
		1:  "one",
		10: "ten",
		3:  "three",
		7:  "seven",
	}

	for k, v := range expected {
		sl.Set(k, v)
	}

	for k, v := range expected {
		val, ok := sl.Get(k)
		if !ok {
			t.Errorf("Get(%d) should return ok=true", k)
			continue
		}
		if val != v {
			t.Errorf("Get(%d): expected %q, got %v", k, v, val)
		}
	}
}

// ======================== Delete Tests ========================

func TestDelete_Exists(t *testing.T) {
	sl := NewInt()
	sl.Set(1, "one")
	sl.Set(2, "two")

	val, ok := sl.Delete(1)
	if !ok {
		t.Error("Delete existing key should return ok=true")
	}
	if val != "one" {
		t.Errorf("expected 'one', got %v", val)
	}
	if sl.size != 1 {
		t.Errorf("expected size=1 after delete, got %d", sl.size)
	}

	// Verify key is actually removed
	_, ok = sl.Get(1)
	if ok {
		t.Error("Get after delete should return ok=false")
	}
}

func TestDelete_NotExists(t *testing.T) {
	sl := NewInt()
	sl.Set(1, "one")

	val, ok := sl.Delete(999)
	if ok {
		t.Error("Delete non-existing key should return ok=false")
	}
	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}
	if sl.size != 1 {
		t.Errorf("expected size=1 unchanged, got %d", sl.size)
	}
}

func TestDelete_EmptyList(t *testing.T) {
	sl := NewInt()

	val, ok := sl.Delete(1)
	if ok {
		t.Error("Delete on empty list should return ok=false")
	}
	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}
}

func TestDelete_AllKeys(t *testing.T) {
	sl := NewInt()
	keys := []int{1, 2, 3, 4, 5}
	for _, k := range keys {
		sl.Set(k, k*10)
	}

	for _, k := range keys {
		val, ok := sl.Delete(k)
		if !ok {
			t.Errorf("Delete(%d) should return ok=true", k)
		}
		if val != k*10 {
			t.Errorf("Delete(%d): expected %d, got %v", k, k*10, val)
		}
	}

	if sl.size != 0 {
		t.Errorf("expected size=0 after deleting all, got %d", sl.size)
	}
}

func TestDelete_ThenReinsert(t *testing.T) {
	sl := NewInt()

	sl.Set(1, "first")
	sl.Delete(1)
	old, existed := sl.Set(1, "second")

	if existed {
		t.Error("re-insert after delete should return existed=false")
	}
	if old != nil {
		t.Errorf("expected old=nil, got %v", old)
	}

	val, ok := sl.Get(1)
	if !ok {
		t.Error("Get re-inserted key should return ok=true")
	}
	if val != "second" {
		t.Errorf("expected 'second', got %v", val)
	}
}

// ======================== randomLevel Tests ========================

func TestRandomLevel_Range(t *testing.T) {
	sl := NewInt()

	for i := 0; i < 10000; i++ {
		lvl := sl.randomLevel()
		if lvl < 1 {
			t.Errorf("randomLevel should be >= 1, got %d", lvl)
		}
		if lvl > sl.maxLevel {
			t.Errorf("randomLevel should be <= maxLevel(%d), got %d", sl.maxLevel, lvl)
		}
	}
}

func TestRandomLevel_Distribution(t *testing.T) {
	sl := NewInt()
	counts := make([]int, sl.maxLevel+1)

	const samples = 100000
	for i := 0; i < samples; i++ {
		lvl := sl.randomLevel()
		counts[lvl]++
	}

	// With P=2, about 50% should be level 1, 25% level 2, 12.5% level 3, etc.
	// Allow loose tolerance for randomness.
	if counts[1] < samples/4 {
		t.Errorf("level 1 count too low: %d / %d", counts[1], samples)
	}
	if counts[1] < counts[2] {
		t.Error("level 1 count should be >= level 2 count")
	}
	if counts[2] < counts[3] && counts[2] > 0 {
		t.Error("level 2 count should be >= level 3 count")
	}
}

// ======================== Integration Tests ========================

func TestInsertGetDelete_Sequential(t *testing.T) {
	sl := NewInt()

	// Insert
	for i := 0; i < 100; i++ {
		sl.Set(i, i*10)
	}
	if sl.size != 100 {
		t.Fatalf("expected size=100, got %d", sl.size)
	}

	// Get all
	for i := 0; i < 100; i++ {
		val, ok := sl.Get(i)
		if !ok {
			t.Errorf("Get(%d) should return ok=true", i)
			continue
		}
		if val != i*10 {
			t.Errorf("Get(%d): expected %d, got %v", i, i*10, val)
		}
	}

	// Delete even keys
	for i := 0; i < 100; i += 2 {
		val, ok := sl.Delete(i)
		if !ok {
			t.Errorf("Delete(%d) should return ok=true", i)
		}
		if val != i*10 {
			t.Errorf("Delete(%d): expected %d, got %v", i, i*10, val)
		}
	}
	if sl.size != 50 {
		t.Errorf("expected size=50 after deleting half, got %d", sl.size)
	}

	// Odd keys still exist
	for i := 1; i < 100; i += 2 {
		_, ok := sl.Get(i)
		if !ok {
			t.Errorf("Get(%d) should still return ok=true", i)
		}
	}

	// Even keys are gone
	for i := 0; i < 100; i += 2 {
		_, ok := sl.Get(i)
		if ok {
			t.Errorf("Get(%d) should return ok=false after delete", i)
		}
	}
}

func TestInsert_ReverseOrder(t *testing.T) {
	sl := NewInt()

	for i := 100; i >= 1; i-- {
		sl.Set(i, fmt.Sprintf("v%d", i))
	}

	if sl.size != 100 {
		t.Fatalf("expected size=100, got %d", sl.size)
	}

	for i := 1; i <= 100; i++ {
		val, ok := sl.Get(i)
		if !ok {
			t.Errorf("Get(%d) should return ok=true", i)
			continue
		}
		expected := fmt.Sprintf("v%d", i)
		if val != expected {
			t.Errorf("Get(%d): expected %q, got %v", i, expected, val)
		}
	}
}

func TestInsert_RandomOrder(t *testing.T) {
	sl := NewInt()
	rng := rand.New(rand.NewSource(42))
	keys := make([]int, 50)
	for i := range keys {
		keys[i] = i
	}
	rng.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	for _, k := range keys {
		sl.Set(k, k*100)
	}

	sort.Ints(keys)
	for _, k := range keys {
		val, ok := sl.Get(k)
		if !ok {
			t.Errorf("Get(%d) should return ok=true", k)
			continue
		}
		if val != k*100 {
			t.Errorf("Get(%d): expected %d, got %v", k, k*100, val)
		}
	}
}

func TestLargeDataset(t *testing.T) {
	sl := NewInt()
	const n = 10000

	for i := 0; i < n; i++ {
		sl.Set(i, i*2)
	}

	if sl.size != n {
		t.Fatalf("expected size=%d, got %d", n, sl.size)
	}

	// Spot-check gets
	for i := 0; i < n; i += 1000 {
		val, ok := sl.Get(i)
		if !ok {
			t.Errorf("Get(%d) should return ok=true", i)
		}
		if val != i*2 {
			t.Errorf("Get(%d): expected %d, got %v", i, i*2, val)
		}
	}

	// Delete half
	for i := 0; i < n; i += 2 {
		sl.Delete(i)
	}

	if sl.size != n/2 {
		t.Errorf("expected size=%d, got %d", n/2, sl.size)
	}
}

// ======================== Edge Cases ========================

func TestStringKeys(t *testing.T) {
	sl := New(func(a, b any) int {
		as, bs := a.(string), b.(string)
		switch {
		case as < bs:
			return -1
		case as > bs:
			return 1
		default:
			return 0
		}
	})

	sl.Set("banana", 2)
	sl.Set("apple", 1)
	sl.Set("cherry", 3)

	val, ok := sl.Get("apple")
	if !ok || val != 1 {
		t.Errorf("Get('apple'): expected 1, got %v (ok=%v)", val, ok)
	}

	val, ok = sl.Get("banana")
	if !ok || val != 2 {
		t.Errorf("Get('banana'): expected 2, got %v (ok=%v)", val, ok)
	}

	val, ok = sl.Get("cherry")
	if !ok || val != 3 {
		t.Errorf("Get('cherry'): expected 3, got %v (ok=%v)", val, ok)
	}

	val, ok = sl.Get("durian")
	if ok {
		t.Errorf("Get('durian') should return ok=false, got %v", val)
	}
}

func TestStructKeys(t *testing.T) {
	type Point struct{ X, Y int }

	sl := New(func(a, b any) int {
		pa, pb := a.(Point), b.(Point)
		if pa.X != pb.X {
			return pa.X - pb.X
		}
		return pa.Y - pb.Y
	})

	sl.Set(Point{0, 0}, "origin")
	sl.Set(Point{1, 2}, "p1")
	sl.Set(Point{2, 4}, "p2")

	val, ok := sl.Get(Point{1, 2})
	if !ok || val != "p1" {
		t.Errorf("Get(Point{1,2}): expected 'p1', got %v (ok=%v)", val, ok)
	}

	val, ok = sl.Get(Point{0, 0})
	if !ok || val != "origin" {
		t.Errorf("Get(Point{0,0}): expected 'origin', got %v (ok=%v)", val, ok)
	}

	_, ok = sl.Get(Point{99, 99})
	if ok {
		t.Error("Get non-existing Point should return ok=false")
	}
}

func TestUpdateValueToNil(t *testing.T) {
	sl := NewInt()

	sl.Set(1, "value")
	old, existed := sl.Set(1, nil)

	if !existed {
		t.Error("update should return existed=true")
	}
	if old != "value" {
		t.Errorf("expected old='value', got %v", old)
	}

	val, ok := sl.Get(1)
	if !ok {
		t.Error("Get should return ok=true for existing key with nil value")
	}
	if val != nil {
		t.Errorf("expected nil value, got %v", val)
	}
}

func TestSameKeyDifferentValues(t *testing.T) {
	sl := NewInt()

	type testCase struct {
		setValue  any
		getValue  any
	}

	cases := []testCase{
		{42, 42},
		{"string", "string"},
		{true, true},
		{3.14, 3.14},
	}

	for _, tc := range cases {
		sl.Set(1, tc.setValue)
		val, ok := sl.Get(1)
		if !ok {
			t.Errorf("Get after Set(%v) should return ok=true", tc.setValue)
		}
		if val != tc.getValue {
			t.Errorf("Get: expected %v, got %v", tc.getValue, val)
		}
	}
}

// ======================== Benchmarks ========================

func BenchmarkSet(b *testing.B) {
	sl := NewInt()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Set(i, i)
	}
}

func BenchmarkGet(b *testing.B) {
	sl := NewInt()
	for i := 0; i < 100000; i++ {
		sl.Set(i, i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Get(i % 100000)
	}
}

func BenchmarkDelete(b *testing.B) {
	sl := NewInt()
	for i := 0; i < b.N; i++ {
		sl.Set(i, i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Delete(i)
	}
}
