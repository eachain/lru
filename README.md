# lru



## 示例

```go
package main

import (
	"fmt"

	"github.com/eachain/lru"
)

func main() {
	cache := lru.New[string, int](3)
	cache.Set("banana", 3)
	cache.Set("apple", 2)
	cache.Set("pear", 4)
	cache.Set("orange", 1)

	cache.All()(func(s string, i int) bool {
		fmt.Printf("%v:%v ", s, i)
		return true
	})
	// Output:
	// orange:1 pear:4 apple:2
}
```

