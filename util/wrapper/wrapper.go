package wrapper

import "github.com/gin-gonic/gin"

type Wrapper struct{}

func NewWrapper() *Wrapper {
	return &Wrapper{}
}

func (w *Wrapper) ErrResp(err error) gin.H {
	return gin.H{"error": err.Error()}
}
