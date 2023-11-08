package skiplist

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

// Global maximum level of skip list.
const MaxLevel int = 30

// Defines the data that will be stored in the list.
type Data struct {
	Key   string
	Value []byte
}

// Compares two data keys in a sorting context. Returns
// 1 if data a is larger than b,
// 0 if data a is the same as b,
// -1 if data a is smaller than b.
func compare(a Data, b Data) int {
	return strings.Compare(a.Key, b.Key)
}

// Defines the levels in the skip list
type Level struct {
	Next *Node
}

// Defines the skip list
type SkipList struct {
	// Max possible level for this skip list.
	MaxLevel int

	// Current level of skip list.
	Level int

	// Pointer to header node.
	Header *Node
}

// Defines a Node in the list
type Node struct {
	// Data is the data stored in a node.
	Data Data

	// Forward is a slice containing nodes in a level
	// that are linked from this node.
	Forward []*Node
}

// Creates a new skip list.
func New(maxLevel int) SkipList {

	// Create new node with dummy data as header.
	node := Node{
		Data: Data{
			Key:   "",
			Value: []byte(""),
		},
		Forward: make([]*Node, maxLevel+1),
	}

	return SkipList{
		MaxLevel: maxLevel,
		Level:    0,
		Header:   &node,
	}
}

// Generates a random integer ranging from 0 to the max level of the skip list.
func (s *SkipList) randomLevel() int {
	return rand.Intn(s.MaxLevel)
}

// Inserts data to the list if key does not exist already.
// If the key already exists, the value will be updated with the new one.
func (s *SkipList) Set(key string, value []byte) {
	data := Data{
		Key:   key,
		Value: value,
	}

	curr := s.Header
	update := make([]*Node, s.MaxLevel+1)

	for i := s.Level; i >= 0; i-- {
		for curr.Forward[i] != nil && compare(curr.Forward[i].Data, data) == -1 {
			curr = curr.Forward[i]
		}

		if curr.Forward[i] != nil && compare(curr.Forward[i].Data, data) == 0 {
			curr.Forward[i].Data = data
		}

		update[i] = curr
	}

	curr = curr.Forward[0]
	if curr != nil {
		fmt.Println(compare(curr.Data, data))
	}

	if curr == nil || compare(curr.Data, data) != 0 {
		rLevel := s.randomLevel()

		if rLevel > s.Level {
			for i := s.Level + 1; i < rLevel+1; i++ {
				update[i] = s.Header
			}

			s.Level = rLevel
		}

		n := Node{
			Data:    data,
			Forward: make([]*Node, rLevel+1),
		}

		for i := 0; i <= rLevel; i++ {
			n.Forward[i] = update[i].Forward[i]
			update[i].Forward[i] = &n
		}
	}
}

// Search data from the list.
func (s *SkipList) Search(key string) ([]byte, error) {

	data := Data{
		Key: key,
	}

	curr := s.Header

	for i := s.Level; i >= 0; i-- {
		for curr.Forward[i] != nil && compare(curr.Forward[i].Data, data) == -1 {
			curr = curr.Forward[i]
		}
	}

	curr = curr.Forward[0]

	if curr != nil && compare(curr.Data, data) == 0 {
		return curr.Data.Value, nil
	}

	return nil, errors.New("data not found")
}

// Deletes a data from the list with specified data.
// The data is compared using the compare() function.
func (s *SkipList) Delete(key string) {

	data := Data{
		Key: key,
	}

	curr := s.Header

	update := make([]*Node, s.MaxLevel+1)

	for i := s.Level; i >= 0; i-- {
		for curr.Forward[i] != nil && compare(curr.Forward[i].Data, data) == -1 {
			curr = curr.Forward[i]
		}
		update[i] = curr
	}

	curr = curr.Forward[0]

	if curr != nil && compare(curr.Data, data) == 0 {
		for i := 0; i <= s.Level; i++ {
			if update[i].Forward[i] != curr {
				break
			}

			update[i].Forward[i] = curr.Forward[i]
		}

		for s.Level > 0 && s.Header.Forward[s.Level] == nil {
			s.Level--
		}
	}

}

// Prints all the elements at the bottom level of the list.
func (s *SkipList) Print() {
	for i := 0; i < s.Level; i++ {
		node := s.Header.Forward[i]
		for node != nil {
			fmt.Print(node.Data.Key, " ", node.Data.Value)
			fmt.Print(" ")
			node = node.Forward[i]
		}

	}
	fmt.Println()
}
