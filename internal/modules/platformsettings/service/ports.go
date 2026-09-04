package service

import "dvra-api/internal/shared/models"

// PlatformSettingsRepository es el puerto de salida hacia la persistencia de la
// configuración global del SaaS.
//
// Se declara en el consumidor (ESTANDARES_DESARROLLO.md §3.3 regla 3). Antes el
// service dependía del struct concreto del repositorio, lo que impedía probarlo
// sin base de datos.
//
// GetCached sirve la copia en memoria (TTL de 5 minutos) para la ruta pública,
// que la consulta el frontend en cada carga; Get va siempre a la base.
type PlatformSettingsRepository interface {
	Get() (*models.PlatformSettings, error)
	GetCached() (*models.PlatformSettings, error)
	Update(settings *models.PlatformSettings) error
}
