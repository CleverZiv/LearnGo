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

## Robust

### Panic

#### 什么是 panic

- Panic 是一个运行时异常
- 可以通过调用`panic()`来手动触发 panic

调用 `panic()`之后会依次发生：

1. 当前函数的执行停止
2. defer 关键字后面的函数会被执行
3. 程序执行终止

#### 什么时候使用 panic

- 程序执行遇到不可恢复的错误
  - 如进行 http 连接时，发现 port 已经被占用
- 人为产生的错误
  - 如函数接收参数为指针，不能为空，但传入了 nil

#### 如何处理 panic

使用`defer`

- 使用 defer 书写 panic 之后的恢复逻辑
- 把 defer 写在目标函数的开头
- 多个 defer 的执行顺序是反序（在前面的 defer 后执行）的

```go
// bad
func main() {
  fmt.Println("Enter function main.")
  // trigger panic
  panic(errors.New("This is a panic!"))
  p := recover()
  fmt.Printf("panic:%s\n",p)
  fmt.Println("Exit function main.")
  
}
// good
func main() {
  fmt.Println("Enter function main.")
  defer func() {
    fmt.Println("Enter defer function.")
    if p := recover(); p != nil {
      fmt.Printf("panic:%s\n",p)
    }
    fmt.Println("Exit defer function.")
  }()
  // trigger panic
  panic(errors.New("This is a panic!"))
  fmt.Println("Exit function main.")
  
}
```

### Error

#### 什么是 error

`error`在 go 中定义为一个接口：

```go
type error interface {
  Error() string
}
```

#### 如何自定义 error

自定义的 `error`需要实现`Error()`和`Unwrap()`两个方法

```go
type QueryError struct{
  Query string
  Err error 
}
func (e *QueryError) Error() string {return e.Query + ": " + e.Err.Error()}
func (e *QueryError) Unwrap() string {return e.Err}
```

#### 如何优雅的处理 error

```go
// bad
func foo() error {
  var err error
  if err != nil {
    return err
  }
  return nil
}

// good 加上上下文环境再返回
func foo() error {
  var err error
  if err != nil {
    return fmt.Errorf("context:%w", err)
  }
  return nil
}
```



#### 校验 error 是否匹配

```go
// bad (Go < 1.13) 
func main() {
  e := io.EOF
  fmt.Println(e == io.EOF) // true
  
  e =  fmt.Println("context:%w", io.EOF)
  fmt.Println(e == io.EOF) // false
}
// good 高版本下可以这样做
func main() {
  e := io.EOF
  fmt.Println(errors.Is(err,io.EOF)) // true
  
  e =  fmt.Errorf("context:%w", io.EOF)
  fmt.Println(errors.Is(err,io.EOF)) // true
}
```

#### error 转换

```go
// bad
func main() {
  _, err := os.Open("non-existing")
  _, ok := err.(*fs.PathError)
  fmt.Println(ok) // true
  
  err2 := fmt.Errorf("context:%w", err)
  _, ok := err.(*fs.PathError)
  fmt.Println(ok) // false
}

// good
func main() {
  _, err := os.Open("non-existing")
  var pathError *fs.PathError
  fmt.Println(errors.As(err, &pathError)) // true
  
  err2 := fmt.Errorf("context:%w", err)
  fmt.Println(errors.As(err2, &pathError)) // true
}

```



## Efficient

### 什么是指针

**指针是一个存储了另外一个变量的地址的变量**

### 什么时候使用指针

good：

> - 方法/函数 修改了它的接收者或参数
> - 大容量数据
> - 所有的方法都使用了指针接收者，为保持代码的一致性

bad

>- 指针指向一个接口
>- 方法/函数 未修改它的接收者或参数
>
>- 小容量的数据

### 技巧

- 不要直接使用循环中的变量地址

```go
// bad
func main() {
  x := map[int]string{1:"a",2:"b",3:"c"}
  r := make([]*string, 0 len(x))
  for _, v := range x {
    r = append(r, &v)
  }
  // 打印出 c,c,c 为什么？go 语言会复用，具体的原理可再深入剖析源码
  for _, v := range r {
    fmt.Println(*v)
  }
}
// good
func main() {
  x := map[int]string{1:"a",2:"b",3:"c"}
  r := make([]*string, 0 len(x))
  for _, v := range x {
    c := v // 赋给一个临时变量，不直接使用
    r = append(r, &c)
  }
  // 打印出 c,c,c
  for _, v := range r {
    fmt.Println(*v)
  }
}
```

## Deprecation

### 什么是 deprecation

方法过时

### 如何将方法置为过时

1. 为过时方法添加全面的单元测试
2. 将新方法和过时方法一同测试，保证新方法与旧方法行为一致

#### 标记过时方法

```go
// LegacyFunc does some....
//
// Deprecated: use CoolFunc instead
func LegacyFuc() {
  // do someing ...
}
```

方法上方添加注释，注意：

1. 单词“Deprecated”开头
2. 提供详细的迁移信息
3. 上一行必须为空行
