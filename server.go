package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/zeyadyasser/autocom/engine"
	"github.com/zeyadyasser/autocom/engine/skip"
)

type data struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

func initHTTPServer(config *Config) {
	opts := skip.Options{
		MaxLevels: config.EngineMaxLevels,
		ToLower:   config.EngineToLower,
		SkipBegin: config.EngineSkipBegin,
	}
	E := skip.NewSkipEngine(opts, nil)

	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	server.Use(middleware.CORS())
	server.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 9}))
	server.Use(middleware.Recover())
	// TODO: enhance Auth
	server.Use(middleware.BasicAuth(func(user, password string, c echo.Context) (bool, error) {
		if len(config.Password) == 0 {
			return true, nil
		}
		if user == config.User && password == config.Password {
			return true, nil
		}
		return false, nil
	}))

	server.POST("/set", set(E))
	server.DELETE("/remove", remove(E))
	server.GET("/topn", topN(E))
	server.Logger.Fatal(server.Start(fmt.Sprintf(":%d", config.Port)))
}

func set(e engine.Engine) echo.HandlerFunc {
	return func (c echo.Context) error {
		KV := new(data)
		if err := c.Bind(KV); err != nil {
			return err
		}
		err := e.Set(KV.Key, KV.Value)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	}
}

func remove(e engine.Engine) echo.HandlerFunc {
	return func (c echo.Context) error {
		key := c.QueryParam("key")
		err := e.Remove(key)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	}
}

func topN(e engine.Engine) echo.HandlerFunc {
	return func (c echo.Context) error {
		key := c.QueryParam("key")
		cnt, err := strconv.Atoi(c.QueryParam("cnt"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		res, err := e.TopN(key, cnt)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, res)
	}
}
