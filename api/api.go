package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NewAPI() API {
	r := gin.Default()
	api := API{
		r: r,
	}
	api.setRoute()
	return api
}

type API struct {
	r *gin.Engine
}

func (api API) short(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"message": "ok"})
}

func (api API) decrypt(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"message": "ok"})
}
func (api API) setRoute() {
	api.r.POST("/short", api.short)
	api.r.GET("/", api.decrypt)
}
func (api API) Run(port int) error {
	address := ":" + strconv.Itoa(port)
	err := api.r.Run(address)
	if err != nil {
		return err
	}
	return nil
}
