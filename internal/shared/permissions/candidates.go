package permissions

// Permisos del módulo Candidates
const (
	CandidatesView         = "candidates.view"
	CandidatesCreate       = "candidates.create"
	CandidatesUpdate       = "candidates.update"
	CandidatesDelete       = "candidates.delete"
	CandidatesUploadResume = "candidates.upload_resume"
)

func init() {
	grant(RoleAdmin, CandidatesView, CandidatesCreate, CandidatesUpdate, CandidatesDelete, CandidatesUploadResume)
	grant(RoleRecruiter, CandidatesView, CandidatesCreate, CandidatesUpdate, CandidatesUploadResume)
	// hiring_manager y user ven candidatos solo de sus jobs (matriz 3.2) —
	// el filtrado por job asignado es a nivel de recurso, pendiente (RN-MEMB-007).
	grant(RoleHiringManager, CandidatesView)
	grant(RoleUser, CandidatesView)
}
