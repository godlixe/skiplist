package skiplist

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
)

// Global maximum level of skip list.
const MaxLevel int = 30

var ErrKeyNotFound = errors.New("key not found")

// Used to iterate the skiplist elements
// in a sorted order.
type Iterator struct {
	current *Node
}

type NewData[T any] interface {

	// Compares two data keys in a sorting context. Returns
	// 1 if data a is larger than b,
	// 0 if data a is the same as b,
	// -1 if data a is smaller than b.
	Compare(a, b *T) bool
}

// Defines the data that will be stored in the list.
type Data struct {
	Key   string
	Value []byte
}

// Compares two data keys in a sorting context. Returns
// 1 if data a is larger than b,
// 0 if data a is the same as b,
// -1 if data a is smaller than b.
func compare(a, b *Data) int {
	return strings.Compare(a.Key, b.Key)
}

// Defines the skip list
type SkipList struct {
	mu sync.RWMutex

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

// Creates a new skip list with the default max level of 30
func NewDefault() SkipList {
	return New(MaxLevel)
}

// Generates a random integer ranging from 0 to the max level of the skip list.
func (s *SkipList) randomLevel() int {
	return rand.Intn(s.MaxLevel)
}

// Inserts data to the list if key does not exist already.
// If the key already exists, the value will be updated with the new one.
func (s *SkipList) Set(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data := Data{
		Key:   key,
		Value: value,
	}

	curr := s.Header
	update := make([]*Node, s.MaxLevel+1)

	for i := s.Level; i >= 0; i-- {
		for curr.Forward[i] != nil && compare(&curr.Forward[i].Data, &data) == -1 {
			curr = curr.Forward[i]
		}

		if curr.Forward[i] != nil && compare(&curr.Forward[i].Data, &data) == 0 {
			curr.Forward[i].Data = data
		}

		update[i] = curr
	}

	curr = curr.Forward[0]

	if curr == nil || compare(&curr.Data, &data) != 0 {
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
func (s *SkipList) Search(key string) (*Data, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := Data{
		Key: key,
	}

	curr := s.Header

	for i := s.Level; i >= 0; i-- {
		for curr.Forward[i] != nil && compare(&curr.Forward[i].Data, &data) == -1 {
			curr = curr.Forward[i]
		}
	}

	curr = curr.Forward[0]

	if curr != nil && compare(&curr.Data, &data) == 0 {
		return &curr.Data, nil
	}

	return nil, ErrKeyNotFound
}

// Deletes a data from the list with specified data.
// The data is compared using the compare() function.
func (s *SkipList) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data := Data{
		Key: key,
	}

	curr := s.Header

	update := make([]*Node, s.MaxLevel+1)

	for i := s.Level; i >= 0; i-- {
		for curr.Forward[i] != nil && compare(&curr.Forward[i].Data, &data) == -1 {
			curr = curr.Forward[i]
		}
		update[i] = curr
	}

	curr = curr.Forward[0]

	if curr != nil && compare(&curr.Data, &data) == 0 {
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
	for _, v := range s.Sorted() {
		fmt.Print(v, " ")
	}
	fmt.Println("")
}

// Returns a slice containing all the elements of the skiplist in sorted order.
func (s *SkipList) Sorted() []Data {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var res []Data

	node := s.Header.Forward[0]
	for node != nil {
		res = append(res, node.Data)
		node = node.Forward[0]
	}

	return res
}

func (s *SkipList) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var len int
	node := s.Header.Forward[0]
	for node != nil {
		len++
		node = node.Forward[0]
	}

	return len
}

func (s *SkipList) Iterate() *Iterator {
	return &Iterator{
		current: s.Header.Forward[0],
	}
}

func (i *Iterator) Valid() bool {
	return i.current != nil
}

func (i *Iterator) Next() {
	if i.Valid() {
		i.current = i.current.Forward[0]
	}
}

func (i *Iterator) Data() Data {
	return i.current.Data
}
