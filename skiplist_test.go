package skiplist

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Data struct {
	Key   string
	Value []byte
}

func cmpData(a, b Data) int {
	if a.Key == b.Key {
		return 0
	} else if a.Key < b.Key {
		return -1
	}

	return 1
}

func TestCompare(t *testing.T) {
	type test struct {
		description string
		a           Data
		b           Data
		want        int
	}

	tests := []test{
		{
			description: "a > b",
			a:           Data{Key: "b"},
			b:           Data{Key: "a"},
			want:        1,
		},
		{
			description: "a > b",
			a:           Data{Key: "b"},
			b:           Data{Key: "abcdefg"},
			want:        1,
		},
		{
			description: "a = b",
			a:           Data{Key: "b"},
			b:           Data{Key: "a"},
			want:        1,
		},
		{
			description: "a < b",
			a:           Data{Key: "a"},
			b:           Data{Key: "b"},
			want:        -1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			got := cmpData(tc.a, tc.b)
			assert.Equal(t, tc.want, got)
		})
	}

}

func TestList(t *testing.T) {
	// create new list
	list := New(10, cmpData)

	// insert("1", []byte("hello"))
	list.Set(Data{
		Key:   "1",
		Value: []byte("1"),
	})
	res, _ := list.Search(Data{Key: "1"})
	assert.Equal(t, []byte("1"), res.Value)

	// insert ("1", []byte("hello"))

	list.Set(Data{Key: "2", Value: []byte{2}})
	list.Set(Data{Key: "3", Value: []byte{3}})
	list.Set(Data{Key: "4", Value: []byte{4}})

	assert.Equal(t, 4, list.Len())

	// test delete
	list.Delete(Data{Key: "2"})

	_, err := list.Search(Data{Key: "2"})

	assert.EqualError(t, err, ErrTargetNotFound.Error())

	list.Print()

	for _, v := range list.Sorted() {
		fmt.Println(v.Key, " ", v.Value)
	}

	list.Set(Data{Key: "1", Value: []byte{0}})
	res, _ = list.Search(Data{Key: "1"})

	assert.Equal(t, []byte{0}, res.Value)

	list.Print()

}

func TestIterator(t *testing.T) {
	list := NewDefault(cmpData)

	list.Set(Data{Key: "a", Value: []byte{1}})
	list.Set(Data{Key: "b", Value: []byte{2}})
	list.Set(Data{Key: "c", Value: []byte{3}})

	it := list.Iterate()

	assert.Equal(t, it.Data().Key, "a")
	assert.Equal(t, it.Data().Value, []byte{1})

	it.Next()

	assert.Equal(t, it.Data().Key, "b")
	assert.Equal(t, it.Data().Value, []byte{2})

	it.Next()

	it.Next()

	assert.Equal(t, it.Valid(), false)
}
