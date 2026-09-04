package domain

// Criterios de consulta del módulo.
//
// Viven en domain y no junto a los DTOs de respuesta porque NO son formato de
// cable: son el criterio con el que el service pide datos al repositorio. Que
// además se rellenen desde la query string es una comodidad del transporte, no
// su razón de ser.
//
// Tenerlos aquí es lo que permite que domain declare los puertos sin importar
// service, que a su vez importa domain (sería un ciclo).

// PlacementFilters acota el listado de colocaciones de una empresa.
type PlacementFilters struct {
	Status           string `form:"status"`
	StaffingClientID *uint  `form:"staffing_client_id"`
	CandidateID      *uint  `form:"candidate_id"`
}

// StaffingClientFilters acota el listado de clientes finales de una empresa.
type StaffingClientFilters struct {
	Status string `form:"status"`
}
