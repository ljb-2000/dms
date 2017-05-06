package context

var context Context

type Context struct {
	Debug   bool
	Address string
}

// Get returns context obj
func Get() *Context {
	return &context
}
