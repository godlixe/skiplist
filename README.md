# skiplist

skiplist is a simple skip list implementation in Go. It stores data in a key value manner. The key is of type string, and the value is a byte array.

A skip list is a probabilistic data structure that allows $O(log\ n)$ average complexity for search and insertion within an ordered sequence of n elements. Fast search is made possible by maintaining a linked hierarchy of subsequences, with each successive subsequence skipping over fewer elements than the previous ones. 

![Skip list.](https://upload.wikimedia.org/wikipedia/commons/thumb/8/86/Skip_list.svg/1024px-Skip_list.svg.png)

In layman's terms, a skip list is a data structure that looks like a linked list and has sorted elements. A node could span to higher levels. Levels are used to make searching faster. Searching is done by checking if the current node's key is smaller. If it is smaller, it will skip it. Else, it will go down a level. Notice that not all node spans to the highest level, so the search can be done by skipping nodes that are not needed to be checked. Hence the name skip list. The generation of a node's level depends on a random number, hence the probabilistic feature.

Try the skip list in your Go program by running

```
go get github.com/godlixe/skiplist
```

and importing `github.com/godlixe/skiplist` on the top of your program.


# Examples

```Go
package main

import (
	"fmt"

	"github.com/godlixe/skiplist"
)

// Data type for the skiplist's data
type Data struct {
	Key   string
	Value []byte
}

// A comparator function to sort the data
func cmpData(a, b Data) int {
	if a.Key == b.Key {
		return 0
	} else if a.Key < b.Key {
		return -1
	}

	return 1
}

func main() {
	// creates a new skiplist with max level of 10 and
	// comparator function cmpData
	list := skiplist.New(10, cmpData)

	// sets the key "a" to store the byte array value of "hi"
	list.Set(Data{
		Key:   "a",
		Value: []byte("hi"),
	})

	// search the list for the key "a"
	res, err := list.Search(Data{
		Key: "a",
	})

	// prints the result
	fmt.Println(res, err)

	// update the key "a" to store the byte array value of "hello"
	list.Set(Data{
		Key:   "a",
		Value: []byte("hello"),
	})

	res, err = list.Search(Data{
		Key: "a",
	})

	fmt.Println(res, err)

	// delete key "a" from the list
	list.Delete(Data{
		Key: "a",
	})

	// searching for a non-existing target will return
	// the error ErrTargetNotFound
	res, err = list.Search(Data{
		Key: "a",
	})
	fmt.Println(res, err)

	list.Set(Data{
		Key:   "b",
		Value: []byte("2"),
	})

	list.Set(Data{
		Key:   "c",
		Value: []byte("3"),
	})

	// Prints the length of the list
	fmt.Println(list.Len())

	// Iterate the list
	for i := list.Iterate(); i.Valid(); i.Next() {
		fmt.Printf("key : %s, value : %s\n", i.Data().Key, i.Data().Value)
	}
}
```