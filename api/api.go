package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"one-shot-url/database"
	"one-shot-url/short"
	"strconv"
)

func NewAPI(short short.Shorter, db database.Interactor) API {
	gin.DefaultWriter = log.Writer()
	r := gin.Default()
	api := API{
		r:          r,
		Shorter:    short,
		Interactor: db,
	}
	api.setRoute()
	return api
}

type API struct {
	r *gin.Engine
	short.Shorter
	database.Interactor
}

func (api API) short(c *gin.Context) {
	req := struct {
		Url string `json:"url"`
	}{
		Url: "",
	}

	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid json."})
		return
	}

	if req.Url == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "you should set url."})
		return
	}

	short := api.Shorter.Generate()
	err = api.Interactor.Store(req.Url, short)
	if err != nil {
		c.JSON(http.StatusBadGateway, map[string]string{"message": "server error."})
		return
	}
	c.JSON(http.StatusOK, map[string]string{"short_url": short})
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
