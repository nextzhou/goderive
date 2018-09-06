# GoSet

生成指定类型的集合操作 go 代码

## 使用方式

1. 在你需要集合操作的类型上方加上注释 `// goset`
2. 在该文件所在目录执行 `goset`, 或者在任意目录下执行 `goset [指定类型所在文件或路径]...`
3. 如果执行成功，将会在相同目录下生成 `gen_goset.go` 文件，其中包含指定类型的集合操作相关实现

> 你可以通过 `-o` 选项来改变生成文件的默认名`gen_goset.go`


例子:

```go
// goset
type YourType struct {/*...*/}
```

执行 `goset <path>`, 将会在同级目录下生成 `gen_goset.go` 文件, 其内容是 `YourType` 集合操作的相关实现:


```go
package xxx

type YourTypeSet map[YourType]struct{}


func NewYourTypeSet (capacity int) YourTypeSet {/*...*/}

func NewYourTypeSetFromSlice(items []YourType) YourTypeSet {/*...*/}

func (set YourTypeSet) Extend(items []YourType) {/*...*/}

func (set YourTypeSet) Len() int {/*...*/}

func (set YourTypeSet) Put(key YourType) {/*...*/}

func (set YourTypeSet) Delete(key YourType) {/*...*/}

func (set YourTypeSet) Contains(key YourType) bool {/*...*/}

func (set YourTypeSet) ContainsAny(keys []YourType) bool {/*...*/}

func (set YourTypeSet) ContainsAll(keys []YourType) bool {/*...*/}
```

## 生成选项

你可以通过在 `// goset` 注释后添加选项来修改生成代码的细节，目前支持以下选项：

### rename

通过 `rename` 选项可以修改生成的集合类型的名称(默认名称为`<YourType>Set`)，例如：

```go
// goset: rename=Files
type File struct {/*...*/}
```

生成的集合集合类型名为:

```go
type Files map[File]struct{}
/*...*/
```

### bypointer

生成代码中，默认是存储指定类型的值类型，通过 `bypointer` 选项使生成代码存储指定类型的指针
