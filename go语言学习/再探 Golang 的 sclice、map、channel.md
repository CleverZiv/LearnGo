# 再探 Golang 的 sclice、map、channel 

## slice vs array

### 第一个问题

- Array 需指明长度，长度为常量且不可改变
- Array 长度为其类型中的组成部分
- Array 在作为函数参数的时候会产生 copy
- Golang 所有函数参数都是值传递

#### 什么时候应该用 array？

### 第二个问题

#### slice 的扩容机制

`growslice`方法：

- 当 cap < 1024 时，每次 *2
- 当 cap >=1024 时，每次 *1.25

请问下述情况下，slice 的 len 和 cap 分别是多少？

```go
var s []int
s = append(s,0,1,2) // 3,3
```

#### append slice 的三种方法对比

```go
// 方法1
func BenchmarkAppend(b *teeting.B){
  b.ReportAllocs()
  for i:=0;i<b.N;i++{
    var s []int
    for j:=0;j<10000;j++{
      s=append(s,j)
    }
  }
}
// 方法2 长度不确定，最大容量确定，10000就是最大容量
func BenchmarkAppendAllocated(b *teeting.B){
  b.ReportAllocs()
  for i:=0;i<b.N;i++{
    s:=make([]int,0,10000)
    for j:=0;j<10000;j++{
      s=append(s,j)
    }
  }
}
// 方法3 长度和容量都确定，性能最好
func BenchmarkAppendIndexed(b *teeting.B){
  b.ReportAllocs()
  for i:=0;i<b.N;i++{
     s:=make([]int,10000)
    for j:=0;j<10000;j++{
      s[j]=j
    }
  }
}
```

- 预先分配内存可以提升性能（减少扩容次数）
- 直接使用 index 赋值而非 append 可以提升性能

### 第三个问题

#### case1

```go
func main() {
  var s []int
  for i:=0;i<3;i++{
    s = append(s,i)
  }
  modifySlice(s)
  fmt.Println(s) // [1024 1 2]
}

func modifySlice(s []int) {
  s[0] = 1024
}
```

虽然 golang 是值传递，但是 slice 是的结构体里实际存储的是一个指针，复制出来的 s 仍然持有同样的指针，所以修改原有的内容，在外部会被影响到

#### case2

````go
func main() {
  var s []int
  for i:=0;i<3;i++{
    s = append(s,i)
  }
  modifySlice(s)
  fmt.Println(s) // [1024 1 2]
}

func modifySlice(s []int) {
  s = append(s, 2048)
  s[0] = 1024
}
````

Golang 是值传递，外面的 s 和里面进行 `append(s,2048)`的 s 不是同一个，append 会把里面 s 的 length 改变，但是不会影响外面 s 的 length。所以，**外面的 s 是看不到函数内的 append 操作的**

#### case3

```go
func main() {
  var s []int
  for i:=0;i<3;i++{
    s = append(s,i)
  }
  modifySlice(s)
  fmt.Println(s) // [0 1 2]
}

func modifySlice(s []int) {
  s = append(s, 2048)
  s = append(s, 4096)
  s[0] = 1024
}
```

为什么这里 `s[0]=1024`没有影响到外面的 s 呢？因为里面的 s 扩容了，扩容时会将底层的数组复制一份，那么也就是，外面的 s 和里面的 s 其底层的数组不是一个了。所以修改 s[0] 不会影响外面的 s[0] 的值

#### case4

```go
func main() {
  var s []int
  for i:=0;i<3;i++{
    s = append(s,i)
  }
  modifySlice(s)
  fmt.Println(s) // [1024 1 2]
}

func modifySlice(s []int) {
  s[0] = 1024
  s = append(s, 2048)
  s = append(s, 4096)
}
```

小结：

- 如果没有发生扩容，修改在原来的内存中
- 如果发生了扩容，修改会在新的内存中

### 第四个问题

#### case1

```go
func main() {
  var s []int
  b,_ := json.Marshal(s)
  fmt.Println(string(b)) // null
}

func main() {
  s := []int{}
  b,_ := json.Marshal(s)
  fmt.Println(string(b)) // []
}
```

- 使用 `[]Type{}`或者`make([]Type)`初始化后，slice 不为 nil
- 使用 `var x[]Type`后，slice 为 nil

### 一个小技巧

```go
func normal(s []int) {
  i := 0
  i += s[0]
  i += s[1]
  i += s[2]
  i += s[3]
  println(i)
}

func bce(s []int) {
  _ = s[3]
  
  i := 0
  i += s[0]
  i += s[1]
  i += s[2]
  i += s[3]
  println(i)
}
```

- 如果能确定访问到的 slice 的长度，可以先执行一次让编译器去做优化（减少 bounds checking elimination）

## Map

### 第一个问题

```go
func TestModifyMap(t *testing.T) {
	m := make(map[int]int)
	modifyMap(m)
	t.Log(m) //map[1:1 2:2]
}

func modifyMap(m map[int]int) {
	m[1]=1
	m[2]=2
}
```

`m := make(map[int]int)`返回的是一个指针，因此修改可以被外界看到



结构：

![image-20210601205008178](/Users/lengzefu/Library/Application Support/typora-user-images/image-20210601205008178.png)

- map 实际上的值是指针，传的参数也是指针
- 修改会影响整个 map

### 第二个问题

```go
m := make(map[int]int)
m[1]=1
fmt.Println(&m[1])
```

可以这么写吗？

- map 的 key value 都不可取地址，随着 map 扩容地址会改变
- map 存的是值，读写会发生 copy，如果值大的话会消耗比较大的性能

### 第三个问题

Map 赋值时会自动扩容，删除时会自动缩容吗？

> 不会

- map 删除 key 不会自动缩容
- https:// github.com/golang/go/issues/20135

#### 课后问题

- 为什么同时读写 map 不行？
- 为什么删除 key 没有释放空间？

#### 学习资料

https://qcrao.com/2019/05/22/dive-into-go-map/



## Channel

### 第一个问题

channel 是有锁还是无锁的？

> 有锁

![image-20210601210537882](/Users/lengzefu/Library/Application Support/typora-user-images/image-20210601210537882.png)



![image-20210601210637913](/Users/lengzefu/Library/Application Support/typora-user-images/image-20210601210637913.png)

- Channel 是有锁的
- Channel 底层是个 ringbuffer
- Channel 调用会触发调度
- 高并发、高性能编程不适合使用 channel
- https://github.com/golang/go/issues/8899

### 第二个问题

```go
make(chan Type)
vs
make(chan Type, len) // 带缓冲区
```

- Buffered channel 会发生两次 copy
  - send goroutine -> buf
  - buf -> receive goroutine
- unbuffered channel 会发生一次 copy
  - send goroutine -> receive goroutine
- unbuffered channel receive 完成后 send 才返回
  - happens before 语义，涉及到 golang 的内存模型

### 第三个问题

这段代码有什么问题？

```go
func main() {
  ch := make(chan int)
  go func(){
    ch <- 1
    close(ch)
  }()
  for {
    select {
      case i := <-ch:
      fmt.Println(i)
      default:
      break
    }
  }
}
```

> 这段代码会一直打印 0

- for + select closed channel 会造成死循环
- select 中的 break 无法跳出 for 循环，只跳出 select
- https://qcrao.com/2019/07/22/dive-into-go-channel

