package template

import (
	"net/http"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	"gopkg.in/gin-gonic/gin.v1"
)

var db config.DBPool

const badstatus = http.StatusBadRequest

func Routes(r *gin.Engine) {
	db = config.Con()
	tmpr := r.Group("/api/v1/template")
	tmpr.Use(utils.AuthSessionMidd)
	tmpr.GET("", GetTemplates)
	tmpr.POST("", CreateTemplate)
	tmpr.GET("/:tpl_id", GetATemplate)
	tmpr.PUT("", UpdateTemplate)
	tmpr.DELETE("/:tpl_id", DeleteTemplate)
	tmpr.POST("/action", CreateActionToTmplate)
	tmpr.PUT("/action", UpdateActionToTmplate)
	tmpr.POST("/clone_tpl", CloneTemplate)

	//simple list for ajax use
	tmpr2 := r.Group("/api/v1/template_simple")
	tmpr.Use(utils.AuthSessionMidd)
	tmpr2.GET("", GetTemplatesSimple)
}
