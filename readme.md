# GoSet

生成指定类型的 go 代码

例子:

```go
// goset
type YourType struct {/*...*/}
```

执行 `goset <path>`, 将会在同级目录下生成 `gen_goset.go` 文件, 其内容是 `YourType` 集合操作的相关实现:


```go
package xxx

type YourTypeSet map[YourType]struct{}

func NewYourTypeSet(capacity int) YourTypeSet {/*...*/}

func (set YourTypeSet) Put(key YourType) {/*...*/}

func (set YourTypeSet) Delete(key YourType) {/*...*/}

func (set YourTypeSet) Contains(key YourType) bool {/*...*/}

```
