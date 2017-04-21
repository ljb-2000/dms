package context

var context Context

type Context struct {
	Debug bool
}

func Get() *Context {
	return &context
}
