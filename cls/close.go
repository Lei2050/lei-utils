package cls

import "sync"

type CloseUtil struct {
	closeChan chan struct{}
	once      sync.Once
	closeCb   []func()
}

func NewCloseUtil() *CloseUtil {
	return &CloseUtil{
		closeChan: make(chan struct{}),
	}
}

func MakeCloseUtil() CloseUtil {
	return CloseUtil{
		closeChan: make(chan struct{}),
	}
}

func (c *CloseUtil) IsClosed() bool {
	select {
	case <-c.closeChan:
		return true
	default:
	}
	return false
}

func (c *CloseUtil) C() <-chan struct{} {
	return c.closeChan
}

func (c *CloseUtil) Close(cb func()) {
	c.once.Do(func() {
		close(c.closeChan)

		if cb != nil {
			cb()
		}

		for _, f := range c.closeCb {
			f()
		}
	})
}

func (c *CloseUtil) RegisterCloseCallback(f func()) {
	if f == nil {
		return
	}
	c.closeCb = append(c.closeCb, f)
}
