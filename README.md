# pycall

file: worksapce/test.py
```python
def test():
    raise TypeError("123")

def swap(x: int, y: int) -> (int, int):
    return y, x
```

file: worksapce/main.go
```go
package main

import (
    "fmt"
    "github.com/kkkunny/pycall"
)

func main() {
	err := pycall.InitializeDefault()
	if err != nil {
		panic(err)
	}
	defer pycall.Finalize()

	test, err := pycall.GetFunction[func() error]("test", "test")
	if err != nil {
		panic(err)
	}
	fmt.Println(test())

	swap, err := pycall.GetFunction[func(int, int) (int, int, error)]("test", "swap")
	if err != nil {
		panic(err)
	}
	fmt.Println(swap(1, 2))
}
```
