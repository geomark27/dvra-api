package service

// DashboardRepository es el puerto de salida hacia las estadísticas agregadas
// de una empresa.
//
// El adaptador consulta jobs, candidates y applications con COUNT agregados,
// cruzando tres dominios. Es una excepción consciente a la frontera de módulos
// (§3.3 regla 2), admisible porque es de SOLO LECTURA, agregada y sin lógica de
// negocio: componerla por puertos convertiría tres contadores en N+1 consultas.
// El aislamiento se mantiene: todas las consultas filtran por company_id.
type DashboardRepository interface {
	GetStats(companyID uint) (*DashboardStatsDTO, error)
}
