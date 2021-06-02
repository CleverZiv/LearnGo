# Go 语言的最佳实践

## Readable

### if,else and happy path

#### 避免 else return

```go
// bad example
func absoluteNumber(x int) int{
  if x >= 0 {
    return x
  } else {
    return -x
  }
}
// good example
func absoluteNumber(x int) int{
  if x >= 0 {
    return x
  }
  return -x
}
```

#### 保持 happy path 不要缩进

> Happy path：当所有事情都正常发生时，代码所执行的路径

也就是说，happy path 的逻辑不要嵌套太深，一个表现就是，发生错误时，先返回

```go
// bad example
func aFunc() error {
  err := doSomething()
  if err == nil {
    err := doAnotherThing()
    if err == nil {
      return nil // happy path
    }
    return err
  }
  return err
}
// good example
func aFunc() error {
  if err := doSomething(); err != nil {
    return err
  }
  if err := doAnotherThing(); err != nil {
    return err
  }
  return nil
}
```

#### 去除没有必要的 else

```go
// bad example
var a int
if flag {
  a = 1
}else {
  a = -1
}
// good example
a :=-1
if flag {
  a = 1
}
```



### init

#### 什么是 init()

- `init()`是一个预定义函数
- `init()`函数中的代码执行的顺序比较早
  - 初始化所有导入的 packages
  - 初始化导入 packages 所有的全局变量
  - 执行 `init()`

#### init() 的缺点

- 使代码可读性变差
  - 可以定义在任意 packages 中，并且在程序启动时执行，并不完全可控
- 当 import 一个 package 时，由于 package 中的`init()`可能会带来未知的副作用
  - 如一些全局变量、环境变量被修改

#### good practice

```go
// bad example
package tiktokclients

func init() {
  tiktokClient = buildNewTiktokClient()
  ...
}

// good example
package tiktokclients

func InitTiktokClient(){
  tiktokClient = buildNewTiktokClient()
}

// server.go
func main() {
  tiktokclient.InitTiktokClient()
}
```

将初始化的放进一个普通的函数，然后通过 `main`保证在程序启动时执行。好处：

- 各种客户端初始化如果需要保证顺序，可以在 `main`中保证
- 测试时，可以直接调用函数进行测试；而如果是使用`init()`初始化，没法测试（因为一启动就自动执行了）

```go
// bad example
type MyVar struct {...}

var _defaultVar myVar

func init() {
  _defaultVar = myVar{...}
}
// good example
type MyVar struct {...}

var _defaultVar createVar()

func createVar() MyVar {
  return MyVar{...}
}
```

- good case 就是避免使用 `init`

#### tips

如果一定要使用 `init`，那么有一些几点建议

1. 确保逻辑是确定的，如果不受测试或生产环境的不同而逻辑不同
2. 避免访问全局变量或环境变量
3. 避免依赖不同 `init`的初始化顺序
4. 避免 I/O 操作

### Comments

### 为什么要写注释

提高代码可读性

### 2 种注释风格

1. `/* */ `
2. `// 

### 什么时候写注释

1. package 都需要一个注释
   1. 小包：package foo 的注释写在文件 foo.go 中
   2. 大包：写在目录下的 doc.go 中国
2. 被导出的函数、方法等
3. 避免明显的注释，即代码无复杂逻辑，代码即表达了注释的意思

#### 格式

1. 使用完整的句子
2. 以 identifier 的名字开头，以句号结尾

```go
// bad
// a request to run a command
type Request struct{}

// write the JSON encoding of req to w
func Encode() {}

// good
// Request represents a request to run a command.
type Request struct{}

// Encode writes the JSON encoding of req to w
func Encode() {}

```

这样要求的原因是，如果需要使用 go doc 去生成注释文档时，可读性更高

### What instead of How

注释的内容关注这段代码“做了什么”，而不是“怎么做”

