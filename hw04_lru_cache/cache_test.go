package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic, overflow", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("first", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("second", 200)
		require.False(t, wasInCache)

		wasInCache = c.Set("third", 300)
		require.False(t, wasInCache)

		wasInCache = c.Set("fourth", 400)
		require.False(t, wasInCache)

		val, ok := c.Get("first")
		require.False(t, ok)
		require.Equal(t, nil, val)

		val, ok = c.Get("second")
		require.True(t, ok)
		require.Equal(t, 200, val)

		val, ok = c.Get("third")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("fourth")
		require.True(t, ok)
		require.Equal(t, 400, val)
	})

	t.Run("purge logic, old values", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("first", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("second", 200)
		require.False(t, wasInCache)

		wasInCache = c.Set("third", 300)
		require.False(t, wasInCache)

		// Use first and second values.
		val, ok := c.Get("first")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("second")
		require.True(t, ok)
		require.Equal(t, 200, val)

		// Third value is now the oldest and will be purged the next time another value is added.
		wasInCache = c.Set("fourth", 400)
		require.False(t, wasInCache)

		val, ok = c.Get("first")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("second")
		require.True(t, ok)
		require.Equal(t, 200, val)

		val, ok = c.Get("third")
		require.False(t, ok)
		require.Equal(t, nil, val)

		val, ok = c.Get("fourth")
		require.True(t, ok)
		require.Equal(t, 400, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
