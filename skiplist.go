package skiplist

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

// Global maximum level of skip list.
const MaxLevel int = 30

var ErrTargetNotFound = errors.New("target not found")

// Compares two data keys in a sorting context. Returns
// 1 if data a is larger than b,
// 0 if data a is the same as b,
// -1 if data a is smaller than b.
type Comparator[T any] func(a, b T) int

// Used to iterate the skiplist elements
// in a sorted order.
type Iterator[T any] struct {
	current *Node[T]
}

// Defines the skip list
type SkipList[T any] struct {
	mu sync.RWMutex

	// Max possible level for this skip list.
	MaxLevel int

	// Current level of skip list.
	Level int

	// Pointer to header node.
	Header *Node[T]

	// Comparator function
	Comparator Comparator[T]
}

// Defines a Node in the list
type Node[T any] struct {
	// Data is the data stored in a node.
	Data T

	// Forward is a slice containing nodes in a level
	// that are linked from this node.
	Forward []*Node[T]
}

// Creates a new skip list.
func New[T any](maxLevel int, cmp Comparator[T]) SkipList[T] {
	var zero T

	// Create new node with dummy data as header.
	node := &Node[T]{
		Data:    zero,
		Forward: make([]*Node[T], maxLevel+1),
	}

	return SkipList[T]{
		MaxLevel:   maxLevel,
		Level:      0,
		Header:     node,
		Comparator: cmp,
	}
}

// Creates a new skip list with the default max level of 30
func NewDefault[T any](cmp Comparator[T]) SkipList[T] {
	return New(MaxLevel, cmp)
}

// Generates a random integer ranging from 0 to the max level of the skip list.
func (s *SkipList[T]) randomLevel() int {
	return rand.Intn(s.MaxLevel)
}

// Inserts data to the list if key does not exist already.
// If the key already exists, the value will be updated with the new one.
func (s *SkipList[T]) Set(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	curr := s.Header
	update := make([]*Node[T], s.MaxLevel+1)

	for i := s.Level; i >= 0; i-- {
		for curr.Forward[i] != nil && s.Comparator(curr.Forward[i].Data, value) == -1 {
			curr = curr.Forward[i]
		}

		if curr.Forward[i] != nil && s.Comparator(curr.Forward[i].Data, value) == 0 {
			curr.Forward[i].Data = value
		}

		update[i] = curr
	}

	curr = curr.Forward[0]

	if curr == nil || s.Comparator(curr.Data, value) != 0 {
		rLevel := s.randomLevel()

		if rLevel > s.Level {
			for i := s.Level + 1; i < rLevel+1; i++ {
				update[i] = s.Header
			}

			s.Level = rLevel
		}

		n := Node[T]{
			Data:    value,
			Forward: make([]*Node[T], rLevel+1),
		}

		for i := 0; i <= rLevel; i++ {
			n.Forward[i] = update[i].Forward[i]
			update[i].Forward[i] = &n
		}
	}
}

// Search data from the list.
func (s *SkipList[T]) Search(target T) (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	curr := s.Header

	for i := s.Level; i >= 0; i-- {
		for curr.Forward[i] != nil && s.Comparator(curr.Forward[i].Data, target) == -1 {
			curr = curr.Forward[i]
		}
	}

	curr = curr.Forward[0]

	if curr != nil && s.Comparator(curr.Data, target) == 0 {
		return curr.Data, nil
	}

	var noop T
	return noop, ErrTargetNotFound
}

// Deletes a data from the list with specified data.
// The data is compared using the compare() function.
func (s *SkipList[T]) Delete(target T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	curr := s.Header

	update := make([]*Node[T], s.MaxLevel+1)

	for i := s.Level; i >= 0; i-- {
		for curr.Forward[i] != nil && s.Comparator(curr.Forward[i].Data, target) == -1 {
			curr = curr.Forward[i]
		}
		update[i] = curr
	}

	curr = curr.Forward[0]

	if curr != nil && s.Comparator(curr.Data, target) == 0 {
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
func (s *SkipList[T]) Print() {
	for _, v := range s.Sorted() {
		fmt.Print(v, " ")
	}
	fmt.Println("")
}

// Returns a slice containing all the elements of the skiplist in sorted order.
func (s *SkipList[T]) Sorted() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var res []T

	node := s.Header.Forward[0]
	for node != nil {
		res = append(res, node.Data)
		node = node.Forward[0]
	}

	return res
}

func (s *SkipList[T]) Len() int {
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

func (s *SkipList[T]) Iterate() *Iterator[T] {
	return &Iterator[T]{
		current: s.Header.Forward[0],
	}
}

func (i *Iterator[T]) Valid() bool {
	return i.current != nil
}

func (i *Iterator[T]) Next() {
	if i.Valid() {
		i.current = i.current.Forward[0]
	}
}

func (i *Iterator[T]) Data() T {
	return i.current.Data
}
