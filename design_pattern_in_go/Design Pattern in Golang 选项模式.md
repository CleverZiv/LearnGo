# Design Pattern in Golang: 选项模式

### 问题引入：

假如有一个 `User` 类，在构造对象时，它的属性可填可不填，如何实现可自由组合多种属性的对象构造呢？

```go
type User struct {
  ID     string
  Name   string
  Age    int
  Email  string
  Phone  string
  Gender string  
}
```



## 思路一：多构造函数

Java 中的解决方案是：函数重载。即建立多个函数，每个函数需要的属性是不一样的，使用时会根据具体传了哪些参数决定调用哪个函数。但是在 Golang 中并不支持函数重载，只能用不同的函数名来应对。

```go
func NewUserDefault(id string, name string) (*User, error) {
  return &User{ID: id, Name: name}, nil
}

func NewUserWithPhone(id string, name string, phone string) (*User, error) {
  return &User{ID: id, Name: name, Phone: phone}, nil
}

func NewUserWithEmail(id string, name string, email string) (*User, error) {
  return &User{ID: id, Name: name, Email: email}, nil
}
```

很丑陋，且每次有新的属性组合时就需要添加新的函数，扩展性很差。

## 思路二：配置化

把所有可选的参数放到一个 Config 的 结构体中

```go
type Config struct {
    Age    int
    Email  string
    Phone  string
    Gender string
}
```

然后把Config放到User

```go
type User struct {
    ID   string
    Name string
    Conf *Config
}
```

于是，我们只需要一个 NewUser() 的函数了，但在使用前需要构造 Config 对象。

```go
func NewUser(id string, name string, conf *Config) (*User, error) {
    //...
}
//Using the default configuratrion
user, _ := NewUser("1", "Ada", nil)

conf := Config{Age:18, Phone: "123456"}
user2, _ := NewUser("2", "Bob", &conf)
```

## 思路三：建造者模式

建造者模式会引入一个 `Builder` 对象，但会使得对象在构造时更简单、更自由

## 思路四：选项模式

在 golang 中经常使用选项模式来完成一些基础的服务配置，应对属性非常多的情况。具体实现方法：

```go
package design_pattern_in_go

type User struct {
	ID     string
	Name   string
	Age    int
	Email  string
	Phone  string
	Gender string
}

// Option 是一个函数，具体来说是一个返回 参数为*User 的函数的函数
type Option func(*User)

// WithAge 返回一个 Option函数，返回的这个函数会将 age 赋值给 User 的 Age 属性
func WithAge(age int) Option {
	return func(u *User) {
		u.Age = age
	}
}

// WithEmail 返回一个 Option函数，返回的这个函数会将 email 赋值给 User 的 Email 属性
func WithEmail(email string) Option {
	return func(u *User) {
		u.Email = email
	}
}

func WithPhone(phone string) Option {
	return func(u *User) {
		u.Phone = phone
	}
}

func WithGender(gender string) Option {
	return func(u *User) {
		u.Gender = gender
	}
}

// NewUser 上面返回的那些 Option 函数什么时候执行呢？当然是构造函数里执行
func NewUser(id string, name string, options ...Option)(*User, error) {
	user := User{
		ID: id,
		Name: name,
	}

	for _, option := range options {
		option(&user)
	}

	return &user, nil
}

```

如何使用呢？使用时需要调用构造函数，以及传入 Option 函数

```go
package design_pattern_in_go

import "testing"

func TestNewUser(t *testing.T) {
	user, err := NewUser("1","da", WithAge(20), WithEmail("100231"))
	if err != nil {
		t.Log(err)
	}
	t.Log(user)
}
```

这个看起来比较整洁和优雅，对外的接口只有一个NewUser。

相比于Builder模式，不需要引入一个Builder对象。

对比配置化的模式，也不需要引入一个新的Config。