# （零）Go 语言的上下文

每种语言都有其自带的一些规则，比如命名、访问控制、常用内置函数等。我把它称之为”上下文“，内容可能很零散，但却很重要

#### 1. 命名

Go 语言下的所有的命名都遵循的规则：

- 一个名字必须以一个字母（Unicode字母）或下划线开头，后面可以跟任意数量的字母、数字或下划线
- 区分大小写
- Go 语言中自带的关键字不能用于命令
- 推荐使用**驼峰式**命名，但一些特定的词汇需避免使用大小写混合的写法，如：HTML，可以被命名为`htmlEscape`或`HTMLEscape`，但不要是`HtmlEscape`

#### 2. 关键字

Go语言中类似if和switch的关键字有25个；关键字不能用于自定义名字，只能在特定语法结构中使用。

```go
break      default       func     interface   select
case       defer         go       map         struct
chan       else          goto     package     switch
const      fallthrough   if       range       type
continue   for           import   return      var
```

#### 3. 内建常量/类型/函数

内部预先定义的常量、类型、函数：

```go
内建常量: true false iota nil

内建类型: int int8 int16 int32 int64
          uint uint8 uint16 uint32 uint64 uintptr
          float32 float64 complex128 complex64
          bool byte rune string error

内建函数: make len cap new append copy close delete
          complex real imag
          panic recover
```

#### 4. 作用域

在讲作用域之前，结合以下例子，先确定几个概念：

```go
package main

import "fmt"

var arr = [...]int{1,2,3,4,5,6}

func traverseArray(arr [6]int) {
  a := 11
  fmt.Println(a)
	for i, v := range arr {
    b := 22
    fmt.Println(b)
		fmt.Println("index:",i ,"value",v)
	}
}

func main() {
	traverseArray(arr)
}
```

Go 语言中的作用域依赖于”词法域“

- 词法域：一个区域，变量生效的范围（即作用域）

词法域分为显式和隐式

- 显示词法域：由花括弧`{}`所包含的一系列语句组成的语句块。语句块内部声明的名字是无法被外部块访问的。这个块决定了内部声明的名字的作用域范围，也就是作用域。
  - `traverseArray`用`{}`包裹的函数体是显式词法域，在其内声明的变量`a`，只能在函数体内使用，函数外部无法访问
  - `for i, v := range arr`用`{}`包裹的循环体也是显示词法域，在其内声明的变量`b`只能在循环体中使用
- 隐式词法域：没有使用 `{}` 包裹的部分
  - 主语句块：所有的代码，这里对应的是内置函数、内置类型等对应的作用域，称为**内置作用域**
    - 如 `int`、`len()` 是可以在任何代码中使用和访问的
  - 包语句块：包括该包中所有的源码（一个包可能会包括一个目录下的多个文件），对应**包级作用域**
    - 这里的全局变量`arr`属于包语句块，`traverseArray`不是，但如果将其函数名的首字母改为大写，即`TraverseArray`表明该方法可以被包内的其它文件访问，也就称为包语句块
  - 文件语句块：包括该文件中的所有源码，对应**文件级作用域**
    - `traverseArray`只能在该源码文件中访问
  - `for、if、switch`等语句，对应于**局部作用域**
    - 这里的`for`循环中声明的变量`i,v`即对应局部作用域，只能在`for`中使用

**四种作用域的理解**

以上的四种作用域，从上往下，范围从大到小，为了表述方便，我这里自己将范围大的作用域称为高层作用域，而范围小的称为低层作用域。

对于作用域，有以下几点总结：

- 低层作用域，可以访问高层作用域
- 同一层级的作用域，是相互隔离的
- 低层作用域里声明的变量，会覆盖高层作用域里声明的变量

在这里要注意一下，不要将作用域和生命周期混为一谈。声明语句的作用域对应的是一个源代码的文本区域；它是一个编译时的属性。

而一个变量的生命周期是指程序运行时变量存在的有效时间段，在此时间区域内它可以被程序的其他部分引用；是一个运行时的概念。





