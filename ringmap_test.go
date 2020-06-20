package ringmap_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/prgsmall/ringmap"
	"github.com/stretchr/testify/assert"
)

var ringMapCapacity = 777

func TestObjectCreation(t *testing.T) {

	t.Run("TestNewRingMap", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		assert.IsType(t, &ringmap.RingMap{}, m)
		assert.Equal(t, ringMapCapacity, m.Capacity())
		assert.EqualValues(t, false, m.IsFull())
	})
}

func TestGet(t *testing.T) {
	t.Run("ReturnsNotOKIfStringKeyDoesntExist", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		_, ok := m.Get("foo")
		assert.False(t, ok)
	})

	t.Run("ReturnsNotOKIfNonStringKeyDoesntExist", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		_, ok := m.Get(123)
		assert.False(t, ok)
	})

	t.Run("ReturnsOKIfKeyExists", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set("foo", "bar")
		_, ok := m.Get("foo")
		assert.True(t, ok)
	})

	t.Run("ReturnsValueForKey", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set("foo", "bar")
		value, _ := m.Get("foo")
		assert.Equal(t, "bar", value)
	})

	t.Run("ReturnsDynamicValueForKey", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set("foo", "baz")
		value, _ := m.Get("foo")
		assert.Equal(t, "baz", value)
	})

	t.Run("KeyDoesntExistOnNonEmptyMap", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set("foo", "baz")
		_, ok := m.Get("bar")
		assert.False(t, ok)
	})

	t.Run("ValueForKeyDoesntExistOnNonEmptyMap", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set("foo", "baz")
		value, _ := m.Get("bar")
		assert.Nil(t, value)
	})
}

func TestSet(t *testing.T) {
	t.Run("ReturnsTrueIfStringKeyIsNew", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		ok := m.Set("foo", "bar")
		assert.True(t, ok)
	})

	t.Run("ReturnsTrueIfNonStringKeyIsNew", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		ok := m.Set(123, "bar")
		assert.True(t, ok)
	})

	t.Run("ValueCanBeNonString", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		ok := m.Set(123, true)
		assert.True(t, ok)
	})

	t.Run("ReturnsFalseIfKeyIsNotNew", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set("foo", "bar")
		ok := m.Set("foo", "bar")
		assert.False(t, ok)
	})

	t.Run("SetThreeDifferentKeys", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set("foo", "bar")
		m.Set("baz", "qux")
		ok := m.Set("quux", "corge")
		assert.True(t, ok)
	})
}

func TestLen(t *testing.T) {
	t.Run("EmptyMapIsZeroLen", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		assert.Equal(t, 0, m.Len())
	})

	t.Run("SingleElementIsLenOne", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set(123, true)
		assert.Equal(t, 1, m.Len())
	})

	t.Run("ThreeElements", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set(1, true)
		m.Set(2, true)
		m.Set(3, true)
		assert.Equal(t, 3, m.Len())
	})

	t.Run("ThreeElementsWithMax", func(t *testing.T) {
		m := ringmap.NewRingMap(3)
		assert.Equal(t, false, m.IsFull())
		m.Set(1, true)
		assert.Equal(t, false, m.IsFull())
		m.Set(2, true)
		assert.Equal(t, false, m.IsFull())
		m.Set(3, true)
		assert.Equal(t, 3, m.Len())
		assert.Equal(t, true, m.IsFull())
		assert.Equal(t, m.Front().Key, 1)

		m.Set(4, true)
		assert.Equal(t, 3, m.Len())
		assert.Equal(t, true, m.IsFull())
		assert.Equal(t, m.Front().Key, 2)
		assert.Equal(t, m.Back().Key, 4)
	})
}

func TestKeys(t *testing.T) {
	t.Run("EmptyMap", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		assert.Empty(t, m.Keys())
	})

	t.Run("OneElement", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set(1, true)
		assert.Equal(t, []interface{}{1}, m.Keys())
	})

	t.Run("RetainsOrder", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		for i := 1; i < 10; i++ {
			m.Set(i, true)
		}
		assert.Equal(t,
			[]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9},
			m.Keys())
	})

	t.Run("ReplacingKeyDoesntChangeOrder", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set("foo", true)
		m.Set("bar", true)
		m.Set("foo", false)
		assert.Equal(t,
			[]interface{}{"foo", "bar"},
			m.Keys())
	})

	t.Run("KeysAfterDelete", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set("foo", true)
		m.Set("bar", true)
		m.Delete("foo")
		assert.Equal(t, []interface{}{"bar"}, m.Keys())
	})
}

func TestDelete(t *testing.T) {
	t.Run("KeyDoesntExistReturnsFalse", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		assert.False(t, m.Delete("foo"))
	})

	t.Run("KeyDoesExist", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set("foo", nil)
		assert.True(t, m.Delete("foo"))
	})

	t.Run("KeyNoLongerExists", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set("foo", nil)
		m.Delete("foo")
		_, exists := m.Get("foo")
		assert.False(t, exists)
	})

	t.Run("KeyDeleteIsIsolated", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set("foo", nil)
		m.Set("bar", nil)
		m.Delete("foo")
		_, exists := m.Get("bar")
		assert.True(t, exists)
	})
}

func TestRingMap_Front(t *testing.T) {
	t.Run("NilOnEmptyMap", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		assert.Nil(t, m.Front())
	})

	t.Run("NilOnEmptyMap", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set(1, true)
		assert.NotNil(t, m.Front())
	})
}

func TestRingMap_Back(t *testing.T) {
	t.Run("NilOnEmptyMap", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		assert.Nil(t, m.Back())
	})

	t.Run("NilOnEmptyMap", func(t *testing.T) {
		m := ringmap.NewRingMap(ringMapCapacity)
		m.Set(1, true)
		assert.NotNil(t, m.Back())
	})
}

func benchmarkMap_Set(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := make(map[int]bool)
		for i := 0; i < b.N*multiplier; i++ {
			m[i] = true
		}
	}
}

func BenchmarkMap_Set(b *testing.B) {
	benchmarkMap_Set(1)(b)
}

func benchmarkRingMap_Set(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := ringmap.NewRingMap(ringMapCapacity)
		for i := 0; i < b.N*multiplier; i++ {
			m.Set(i, true)
		}
	}
}

func BenchmarkRingMap_Set(b *testing.B) {
	benchmarkRingMap_Set(1)(b)
}

func benchmarkMap_Get(multiplier int) func(b *testing.B) {
	m := make(map[int]bool)
	for i := 0; i < 1000*multiplier; i++ {
		m[i] = true
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = m[i%1000*multiplier]
		}
	}
}

func BenchmarkMap_Get(b *testing.B) {
	benchmarkMap_Get(1)(b)
}

func benchmarkRingMap_Get(multiplier int) func(b *testing.B) {
	m := ringmap.NewRingMap(ringMapCapacity)
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.Get(i % 1000 * multiplier)
		}
	}
}

func BenchmarkRingMap_Get(b *testing.B) {
	benchmarkRingMap_Get(1)(b)
}

func benchmarkRingMap_Len(multiplier int) func(b *testing.B) {
	m := ringmap.NewRingMap(ringMapCapacity)
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		var temp int
		for i := 0; i < b.N; i++ {
			temp = m.Len()
		}

		// prevent compiler from optimising Len away.
		tempInt = temp
	}
}

func BenchmarkRingMap_Len(b *testing.B) {
	benchmarkRingMap_Len(1)(b)
}

func benchmarkMap_Delete(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := make(map[int]bool)
		for i := 0; i < b.N*multiplier; i++ {
			m[i] = true
		}

		for i := 0; i < b.N; i++ {
			delete(m, i)
		}
	}
}

func BenchmarkMap_Delete(b *testing.B) {
	benchmarkMap_Delete(1)(b)
}

func benchmarkRingMap_Delete(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := ringmap.NewRingMap(ringMapCapacity)
		for i := 0; i < b.N*multiplier; i++ {
			m.Set(i, true)
		}

		for i := 0; i < b.N; i++ {
			m.Delete(i)
		}
	}
}

func BenchmarkRingMap_Delete(b *testing.B) {
	benchmarkRingMap_Delete(1)(b)
}

func benchmarkMap_Iterate(multiplier int) func(b *testing.B) {
	m := make(map[int]bool)
	for i := 0; i < 1000*multiplier; i++ {
		m[i] = true
	}
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range m {
				nothing(v)
			}
		}
	}
}
func BenchmarkMap_Iterate(b *testing.B) {
	benchmarkMap_Iterate(1)(b)
}

func benchmarkRingMap_Iterate(multiplier int) func(b *testing.B) {
	m := ringmap.NewRingMap(ringMapCapacity)
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, key := range m.Keys() {
				_, v := m.Get(key)
				nothing(v)
			}
		}
	}
}

func BenchmarkRingMap_Iterate(b *testing.B) {
	benchmarkRingMap_Iterate(1)(b)
}

func benchmarkRingMap_Keys(multiplier int) func(b *testing.B) {
	m := ringmap.NewRingMap(ringMapCapacity)
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.Keys()
		}
	}
}

func benchmarkMapString_Set(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := make(map[string]bool)
		a := "12345678"
		for i := 0; i < b.N*multiplier; i++ {
			m[a+strconv.Itoa(i)] = true
		}
	}
}

func BenchmarkMapString_Set(b *testing.B) {
	benchmarkMapString_Set(1)(b)
}

func benchmarkRingMapString_Set(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := ringmap.NewRingMap(ringMapCapacity)
		a := "12345678"
		for i := 0; i < b.N*multiplier; i++ {
			m.Set(a+strconv.Itoa(i), true)
		}
	}
}

func BenchmarkRingMapString_Set(b *testing.B) {
	benchmarkRingMapString_Set(1)(b)
}

func benchmarkMapString_Get(multiplier int) func(b *testing.B) {
	m := make(map[string]bool)
	a := "12345678"
	for i := 0; i < 1000*multiplier; i++ {
		m[a+strconv.Itoa(i)] = true
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = m[a+strconv.Itoa(i%1000*multiplier)]
		}
	}
}

func BenchmarkMapString_Get(b *testing.B) {
	benchmarkMapString_Get(1)(b)
}

func benchmarkRingMapString_Get(multiplier int) func(b *testing.B) {
	m := ringmap.NewRingMap(ringMapCapacity)
	a := "12345678"
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(a+strconv.Itoa(i), true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.Get(a + strconv.Itoa(i%1000*multiplier))
		}
	}
}

func BenchmarkRingMapString_Get(b *testing.B) {
	benchmarkRingMapString_Get(1)(b)
}

func benchmarkMapString_Delete(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := make(map[string]bool)
		a := "12345678"
		for i := 0; i < b.N*multiplier; i++ {
			m[a+strconv.Itoa(i)] = true
		}

		for i := 0; i < b.N; i++ {
			delete(m, a+strconv.Itoa(i))
		}
	}
}

func BenchmarkMapString_Delete(b *testing.B) {
	benchmarkMapString_Delete(1)(b)
}

func benchmarkRingMapString_Delete(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := ringmap.NewRingMap(ringMapCapacity)
		a := "12345678"
		for i := 0; i < b.N*multiplier; i++ {
			m.Set(a+strconv.Itoa(i), true)
		}

		for i := 0; i < b.N; i++ {
			m.Delete(a + strconv.Itoa(i))
		}
	}
}

func BenchmarkRingMapString_Delete(b *testing.B) {
	benchmarkRingMapString_Delete(1)(b)
}

func benchmarkMapString_Iterate(multiplier int) func(b *testing.B) {
	m := make(map[string]bool)
	a := "12345678"
	for i := 0; i < 1000*multiplier; i++ {
		m[a+strconv.Itoa(i)] = true
	}
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range m {
				nothing(v)
			}
		}
	}
}
func BenchmarkMapString_Iterate(b *testing.B) {
	benchmarkMapString_Iterate(1)(b)
}

func benchmarkRingMapString_Iterate(multiplier int) func(b *testing.B) {
	m := ringmap.NewRingMap(ringMapCapacity)
	a := "12345678"
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(a+strconv.Itoa(i), true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, key := range m.Keys() {
				_, v := m.Get(key)
				nothing(v)
			}
		}
	}
}

func BenchmarkRingMapString_Iterate(b *testing.B) {
	benchmarkRingMapString_Iterate(1)(b)
}

func BenchmarkRingMap_Keys(b *testing.B) {
	benchmarkRingMap_Keys(1)(b)
}

func ExampleNewRingMap() {
	m := ringmap.NewRingMap(ringMapCapacity)

	m.Set("foo", "bar")
	m.Set("qux", 1.23)
	m.Set(123, true)

	m.Delete("qux")

	for _, key := range m.Keys() {
		value, _ := m.Get(key)
		fmt.Println(key, value)
	}
}

func ExampleRingMap_Front() {
	m := ringmap.NewRingMap(ringMapCapacity)
	m.Set(1, true)
	m.Set(2, true)

	for el := m.Front(); el != nil; el = el.Next() {
		fmt.Println(el)
	}
}

func nothing(v interface{}) {
	v = false
}

func benchmarkBigMap_Set() func(b *testing.B) {
	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			m := make(map[int]bool)
			for i := 0; i < 10000000; i++ {
				m[i] = true
			}
		}
	}
}

func BenchmarkBigMap_Set(b *testing.B) {
	benchmarkBigMap_Set()(b)
}

func benchmarkBigRingMap_Set() func(b *testing.B) {
	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			m := ringmap.NewRingMap(ringMapCapacity)
			for i := 0; i < 10000000; i++ {
				m.Set(i, true)
			}
		}
	}
}

func BenchmarkBigRingMap_Set(b *testing.B) {
	benchmarkBigRingMap_Set()(b)
}

func benchmarkBigMap_Get() func(b *testing.B) {
	m := make(map[int]bool)
	for i := 0; i < 10000000; i++ {
		m[i] = true
	}

	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			for i := 0; i < 10000000; i++ {
				_ = m[i]
			}
		}
	}
}

func BenchmarkBigMap_Get(b *testing.B) {
	benchmarkBigMap_Get()(b)
}

func benchmarkBigRingMap_Get() func(b *testing.B) {
	m := ringmap.NewRingMap(ringMapCapacity)
	for i := 0; i < 10000000; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			for i := 0; i < 10000000; i++ {
				m.Get(i)
			}
		}
	}
}

func BenchmarkBigRingMap_Get(b *testing.B) {
	benchmarkBigRingMap_Get()(b)
}

func benchmarkBigMap_Iterate() func(b *testing.B) {
	m := make(map[int]bool)
	for i := 0; i < 10000000; i++ {
		m[i] = true
	}
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range m {
				nothing(v)
			}
		}
	}
}
func BenchmarkBigMap_Iterate(b *testing.B) {
	benchmarkBigMap_Iterate()(b)
}

func benchmarkBigRingMap_Iterate() func(b *testing.B) {
	m := ringmap.NewRingMap(ringMapCapacity)
	for i := 0; i < 10000000; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, key := range m.Keys() {
				_, v := m.Get(key)
				nothing(v)
			}
		}
	}
}

func BenchmarkBigRingMap_Iterate(b *testing.B) {
	benchmarkBigRingMap_Iterate()(b)
}

func benchmarkBigMapString_Set() func(b *testing.B) {
	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			m := make(map[string]bool)
			a := "1234567"
			for i := 0; i < 10000000; i++ {
				m[a+strconv.Itoa(i)] = true
			}
		}
	}
}

func BenchmarkBigMapString_Set(b *testing.B) {
	benchmarkBigMapString_Set()(b)
}

func benchmarkBigRingMapString_Set() func(b *testing.B) {
	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			m := ringmap.NewRingMap(ringMapCapacity)
			a := "1234567"
			for i := 0; i < 10000000; i++ {
				m.Set(a+strconv.Itoa(i), true)
			}
		}
	}
}

func BenchmarkBigRingMapString_Set(b *testing.B) {
	benchmarkBigRingMapString_Set()(b)
}

func benchmarkBigMapString_Get() func(b *testing.B) {
	m := make(map[string]bool)
	a := "1234567"
	for i := 0; i < 10000000; i++ {
		m[a+strconv.Itoa(i)] = true
	}

	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			for i := 0; i < 10000000; i++ {
				_ = m[a+strconv.Itoa(i)]
			}
		}
	}
}

func BenchmarkBigMapString_Get(b *testing.B) {
	benchmarkBigMapString_Get()(b)
}

func benchmarkBigRingMapString_Get() func(b *testing.B) {
	m := ringmap.NewRingMap(ringMapCapacity)
	a := "1234567"
	for i := 0; i < 10000000; i++ {
		m.Set(a+strconv.Itoa(i), true)
	}

	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			for i := 0; i < 10000000; i++ {
				m.Get(a + strconv.Itoa(i))
			}
		}
	}
}

func BenchmarkBigRingMapString_Get(b *testing.B) {
	benchmarkBigRingMapString_Get()(b)
}

func benchmarkBigMapString_Iterate() func(b *testing.B) {
	m := make(map[string]bool)
	a := "12345678"
	for i := 0; i < 10000000; i++ {
		m[a+strconv.Itoa(i)] = true
	}
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range m {
				nothing(v)
			}
		}
	}
}
func BenchmarkBigMapString_Iterate(b *testing.B) {
	benchmarkBigMapString_Iterate()(b)
}

func benchmarkBigRingMapString_Iterate() func(b *testing.B) {
	m := ringmap.NewRingMap(ringMapCapacity)
	a := "12345678"
	for i := 0; i < 10000000; i++ {
		m.Set(a+strconv.Itoa(i), true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, key := range m.Keys() {
				_, v := m.Get(key)
				nothing(v)
			}
		}
	}
}

func BenchmarkBigRingMapString_Iterate(b *testing.B) {
	benchmarkBigRingMapString_Iterate()(b)
}

func BenchmarkAll(b *testing.B) {
	b.Run("BenchmarkRingMap_Keys", BenchmarkRingMap_Keys)

	b.Run("BenchmarkRingMap_Set", BenchmarkRingMap_Set)
	b.Run("BenchmarkMap_Set", BenchmarkMap_Set)
	b.Run("BenchmarkRingMap_Get", BenchmarkRingMap_Get)
	b.Run("BenchmarkMap_Get", BenchmarkMap_Get)
	b.Run("BenchmarkRingMap_Delete", BenchmarkRingMap_Delete)
	b.Run("BenchmarkMap_Delete", BenchmarkMap_Delete)
	b.Run("BenchmarkRingMap_Iterate", BenchmarkRingMap_Iterate)
	b.Run("BenchmarkMap_Iterate", BenchmarkMap_Iterate)

	b.Run("BenchmarkBigMap_Set", BenchmarkBigMap_Set)
	b.Run("BenchmarkBigRingMap_Set", BenchmarkBigRingMap_Set)
	b.Run("BenchmarkBigMap_Get", BenchmarkBigMap_Get)
	b.Run("BenchmarkBigRingMap_Get", BenchmarkBigRingMap_Get)
	b.Run("BenchmarkBigRingMap_Iterate", BenchmarkBigRingMap_Iterate)
	b.Run("BenchmarkBigMap_Iterate", BenchmarkBigMap_Iterate)

	b.Run("BenchmarkRingMapString_Set", BenchmarkRingMapString_Set)
	b.Run("BenchmarkMapString_Set", BenchmarkMapString_Set)
	b.Run("BenchmarkRingMapString_Get", BenchmarkRingMapString_Get)
	b.Run("BenchmarkMapString_Get", BenchmarkMapString_Get)
	b.Run("BenchmarkRingMapString_Delete", BenchmarkRingMapString_Delete)
	b.Run("BenchmarkMapString_Delete", BenchmarkMapString_Delete)
	b.Run("BenchmarkRingMapString_Iterate", BenchmarkRingMapString_Iterate)
	b.Run("BenchmarkMapString_Iterate", BenchmarkMapString_Iterate)

	b.Run("BenchmarkBigMapString_Set", BenchmarkBigMapString_Set)
	b.Run("BenchmarkBigRingMapString_Set", BenchmarkBigRingMapString_Set)
	b.Run("BenchmarkBigMapString_Get", BenchmarkBigMapString_Get)
	b.Run("BenchmarkBigRingMapString_Get", BenchmarkBigRingMapString_Get)
	b.Run("BenchmarkBigRingMapString_Iterate", BenchmarkBigRingMapString_Iterate)
	b.Run("BenchmarkBigMapString_Iterate", BenchmarkBigMapString_Iterate)
}
