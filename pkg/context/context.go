package context

var context Context

type Context struct {
	Debug bool
}

// Get returns context obj
func Get() *Context {
	return &context
}
