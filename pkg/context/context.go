package context

var context Context

type Context struct {
	Mock    bool
	Debug   bool
	Address string
}

// Get returns context obj
func Get() *Context {
	return &context
}
