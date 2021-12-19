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
func NewUser(id string, name string, options ...Option) (*User, error) {
	user := User{
		ID:   id,
		Name: name,
	}

	for _, option := range options {
		option(&user)
	}

	return &user, nil
}
