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
	define(CandidatesView, CandidatesCreate, CandidatesUpdate, CandidatesDelete, CandidatesUploadResume)
	grant(RoleAdmin, CandidatesView, CandidatesCreate, CandidatesUpdate, CandidatesDelete, CandidatesUploadResume)
	grant(RoleRecruiter, CandidatesView, CandidatesCreate, CandidatesUpdate, CandidatesUploadResume)
	// El repositorio limita al hiring_manager a candidatos vinculados con sus jobs.
	grant(RoleHiringManager, CandidatesView)
}
