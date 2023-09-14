package main

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

func main() {
	engine := gin.New()
	engine.Use(gin.Recovery(), gin.Logger())
	engine.GET("/config/merge", func(ctx *gin.Context) {
		us := ctx.QueryArray("source")
		configs := make([][]byte, 0, len(us))
		for _, u := range us {
			resp, err := http.Get(u)
			if err != nil {
				ctx.AbortWithError(500, err)
				return
			}
			defer resp.Body.Close()
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				ctx.AbortWithError(500, err)
				return
			}
			configs = append(configs, bodyBytes)
		}
		if len(configs) == 0 {
			ctx.String(200, "")
			return
		}
		var others []interface{}
		for _, cfg := range configs[1:] {
			others = append(others, cfg)
		}
		data, err := ini.ShadowLoad(configs[0], others...)
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}
		s, err := ini2String(data)
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}
		ctx.String(200, s)
	})
	engine.Run("0.0.0.0:8080")
}

func ini2String(data *ini.File) (string, error) {
	var buf bytes.Buffer
	_, err := data.WriteTo(&buf)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(buf.String(), `"""`, ""), nil
}
