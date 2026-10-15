package handlers

import (
	"net/http"

	"github.com/bestruirui/bestsub/internal/core/update"
	"github.com/bestruirui/bestsub/internal/server/middleware"
	"github.com/bestruirui/bestsub/internal/server/resp"
	"github.com/bestruirui/bestsub/internal/server/router"
	"github.com/gin-gonic/gin"
)

func init() {
	router.NewGroupRouter("/api/v1/update").
		Use(middleware.Auth()).
		AddRoute(router.NewRoute("", router.GET).Handle(latest))
}

// latest returns the latest BestSub release information.
func latest(c *gin.Context) {
	bestsub, err := update.GetLatestBestsubInfo()
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp.Success(c, map[string]update.LatestInfo{"bestsub": *bestsub})
}
