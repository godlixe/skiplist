package skiplist

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			got := compare(&tc.a, &tc.b)
			assert.Equal(t, tc.want, got)
		})
	}

}

func TestList(t *testing.T) {
	// create new list
	list := New(10)

	// insert("1", []byte("hello"))
	list.Set("1", []byte("1"))
	res, _ := list.Search("1")
	assert.Equal(t, []byte("1"), res)

	// insert ("1", []byte("hello"))

	list.Set("2", []byte("2"))
	list.Set("3", []byte("3"))
	list.Set("4", []byte("4"))

	assert.Equal(t, 4, list.Len())

	// test delete
	list.Delete("2")

	_, err := list.Search("2")
	fmt.Println(err)

	assert.EqualError(t, err, ErrKeyNotFound.Error())

	list.Print()

	for _, v := range list.Sorted() {
		fmt.Println(v.Key, " ", v.Value)
	}

	list.Set("1", []byte("0"))
	res, _ = list.Search("1")
	assert.Equal(t, []byte("0"), res)

	list.Print()

}
