# skiplist

skiplist is a simple skip list implementation in Go.

A skip list is a probabilistic data structure that allows $O(log\ n)$ average complexity for search and insertion within an ordered sequence of n elements. Fast search is made possible by maintaining a linked hierarchy of subsequences, with each successive subsequence skipping over fewer elements than the previous ones. 

![Skip list.](https://upload.wikimedia.org/wikipedia/commons/thumb/8/86/Skip_list.svg/1024px-Skip_list.svg.png)

In layman's terms, a skip list is a data structure that looks like a linked list and has sorted elements. A node could span to higher levels. Levels are used to make searching faster. Searching is done by checking if the current node's key is smaller. If it is smaller, it will skip it. Notice that not all node spans to the highest level, so the search can be done by skipping nodes that are not needed to be checked. Hence the name skip list. The generation of a node's level depends on a random number, hence the probabilistic feature.
