package context

var context Context

type Context struct {
	Debug bool
}

// return context obj
func Get() *Context {
	return &context
}
