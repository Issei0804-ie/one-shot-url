package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"one-shot-url/database/rdb"
	"one-shot-url/short"
	"os"
	"strconv"
)

func NewAPI(short short.Shorter, db rdb.Interactor, port int) API {
	gin.DefaultWriter = log.Writer()
	r := gin.Default()

	domain := os.Getenv("DOMAIN")
	protocol := os.Getenv("PROTOCOL")
	if domain == "" {
		log.Fatal("can not get a domain that is environment value. Did you modify .env?")
	}
	if protocol == "" {
		log.Fatal("can not get a domain that is environment value. Did you modify .env?")
	}
	if protocol != "https" && protocol != "http" {
		log.Fatalf("the protocol(%v)is not supported.\n", protocol)
	}
	isDefaultPort := (protocol == "http" && port == 80) || (protocol == "https" && port == 443)
	api := API{
		r:             r,
		Shorter:       short,
		Interactor:    db,
		Domain:        domain,
		Protocol:      protocol,
		Port:          port,
		IsDefaultPort: isDefaultPort,
	}
	api.setRoute()
	return api
}

type API struct {
	r *gin.Engine
	short.Shorter
	rdb.Interactor
	Domain        string
	Protocol      string
	Port          int
	IsDefaultPort bool
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

	code := api.Shorter.Generate()
	err = api.Interactor.Store(req.Url, code)
	if err != nil {
		c.JSON(http.StatusBadGateway, map[string]string{"message": "server error."})
		return
	}
	domainAndPort := fmt.Sprintf("%v:%v", api.Domain, api.Port)
	if api.IsDefaultPort {
		domainAndPort = api.Domain
	}
	shortURL := fmt.Sprintf("%v://%v/%v", api.Protocol, domainAndPort, code)
	c.JSON(http.StatusOK, map[string]string{"short_url": shortURL})
}

func (api API) decrypt(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "you should set code.\n ex) http://localhost/(code)"})
		return
	}

	longURL, err := api.Interactor.SearchLongURL(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "this URL has not been generated."})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"message": longURL})
}
func (api API) setRoute() {
	api.r.POST("/short", api.short)
	api.r.GET("/:code", api.decrypt)
}
func (api API) Run(port int) error {
	address := ":" + strconv.Itoa(port)
	err := api.r.Run(address)
	if err != nil {
		return err
	}
	return nil
}
