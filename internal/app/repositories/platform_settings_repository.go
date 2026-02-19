package repositories

import (
	"sync"
	"time"

	"dvra-api/internal/app/models"

	"gorm.io/gorm"
)

// PlatformSettingsRepository gestiona la configuración global de la plataforma
type PlatformSettingsRepository struct {
	db *gorm.DB
	// In-memory cache
	cache      *models.PlatformSettings
	cacheMutex sync.RWMutex
	cacheTime  time.Time
	cacheTTL   time.Duration
}

// NewPlatformSettingsRepository crea una nueva instancia del repositorio
func NewPlatformSettingsRepository(db *gorm.DB) *PlatformSettingsRepository {
	return &PlatformSettingsRepository{
		db:       db,
		cacheTTL: 5 * time.Minute, // Cache por 5 minutos
	}
}

// Get obtiene la configuración de la plataforma (siempre de DB)
func (r *PlatformSettingsRepository) Get() (*models.PlatformSettings, error) {
	var settings models.PlatformSettings

	// Intentar obtener el registro existente
	result := r.db.First(&settings)
	if result.Error != nil {
		// Si no existe, crear uno con valores por defecto
		if result.RowsAffected == 0 {
			settings = models.PlatformSettings{
				PlatformName:     "DVRA ATS",
				PlatformShort:    "DVRA",
				Tagline:          "Applicant Tracking System",
				PrimaryColor:     "#2563eb",
				SupportEmail:     "support@dvra.app",
				MarketingURL:     "https://dvra.app",
				DefaultTrialDays: 14,
				DefaultPlanTier:  "free",
			}
			if err := r.db.Create(&settings).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, result.Error
		}
	}

	return &settings, nil
}

// GetCached obtiene la configuración usando cache
func (r *PlatformSettingsRepository) GetCached() (*models.PlatformSettings, error) {
	r.cacheMutex.RLock()

	// Verificar si el cache es válido
	if r.cache != nil && time.Since(r.cacheTime) < r.cacheTTL {
		defer r.cacheMutex.RUnlock()
		return r.cache, nil
	}
	r.cacheMutex.RUnlock()

	// Cache expirado o vacío, obtener de DB
	r.cacheMutex.Lock()
	defer r.cacheMutex.Unlock()

	// Double-check después de obtener el lock
	if r.cache != nil && time.Since(r.cacheTime) < r.cacheTTL {
		return r.cache, nil
	}

	settings, err := r.Get()
	if err != nil {
		return nil, err
	}

	// Actualizar cache
	r.cache = settings
	r.cacheTime = time.Now()

	return settings, nil
}

// Update actualiza la configuración
func (r *PlatformSettingsRepository) Update(settings *models.PlatformSettings) error {
	settings.UpdatedAt = time.Now()

	if err := r.db.Save(settings).Error; err != nil {
		return err
	}

	// Invalidar cache
	r.InvalidateCache()

	return nil
}

// InvalidateCache invalida el cache para forzar recarga
func (r *PlatformSettingsRepository) InvalidateCache() {
	r.cacheMutex.Lock()
	defer r.cacheMutex.Unlock()

	r.cache = nil
	r.cacheTime = time.Time{}
}
