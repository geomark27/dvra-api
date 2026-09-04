package transport

import (
	"dvra-api/internal/modules/publicapi/service"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes monta la career page pública: sin sesión, solo lectura salvo
// la propia postulación.
//
// El aislamiento no viene de un token sino del slug de la empresa en la ruta:
// el service resuelve la empresa y acota todo lo demás a ella.
func RegisterRoutes(public *gin.RouterGroup, svc service.PublicService) {
	h := NewPublicHandler(svc)

	group := public.Group("/public")
	{
		group.GET("/companies/:slug", h.GetCompanyBySlug)
		group.GET("/companies/:slug/jobs", h.GetPublishedJobsByCompany)
		group.GET("/jobs/:id", h.GetPublishedJobByID)
		group.POST("/jobs/:id/apply", h.ApplyToJob)
	}
}
