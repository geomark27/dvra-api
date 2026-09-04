package transport

import (
	"dvra-api/internal/modules/recruitment/service"
	"dvra-api/internal/shared/middleware"
	"dvra-api/internal/shared/permissions"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes monta el núcleo del ATS: departamentos, vacantes, candidatos y
// el pipeline de postulaciones.
//
// Todo vive en el grupo protegido y con scoping por empresa: son los datos de
// tenant por excelencia. Las vacantes también se sirven sin sesión, pero eso lo
// hace el módulo publicapi, que consume estos datos por puerto.
func RegisterRoutes(
	protected *gin.RouterGroup,
	departments service.DepartmentService,
	jobs service.JobService,
	candidates service.CandidateService,
	applications service.ApplicationService,
	rbac middleware.PermissionChecker,
) {
	departmentH := NewDepartmentHandler(departments)
	jobH := NewJobHandler(jobs)
	candidateH := NewCandidateHandler(candidates)
	applicationH := NewApplicationHandler(applications)

	departmentsGroup := protected.Group("/departments")
	{
		departmentsGroup.GET("", middleware.RequirePermission(rbac, permissions.DepartmentsView), departmentH.ListDepartments)
		departmentsGroup.POST("", middleware.RequirePermission(rbac, permissions.DepartmentsCreate), departmentH.CreateDepartment)
		departmentsGroup.GET("/tree", middleware.RequirePermission(rbac, permissions.DepartmentsView), departmentH.GetDepartmentTree)
		departmentsGroup.GET("/:id", middleware.RequirePermission(rbac, permissions.DepartmentsView), departmentH.GetDepartment)
		departmentsGroup.PUT("/:id", middleware.RequirePermission(rbac, permissions.DepartmentsUpdate), departmentH.UpdateDepartment)
		departmentsGroup.DELETE("/:id", middleware.RequirePermission(rbac, permissions.DepartmentsDelete), departmentH.DeactivateDepartment)
	}

	jobsGroup := protected.Group("/jobs")
	{
		jobsGroup.GET("", middleware.RequirePermission(rbac, permissions.JobsView), jobH.GetJobs)
		jobsGroup.POST("", middleware.RequirePermission(rbac, permissions.JobsCreate), jobH.CreateJob)
		jobsGroup.GET("/:id", middleware.RequirePermission(rbac, permissions.JobsView), jobH.GetJob)
		jobsGroup.PUT("/:id", middleware.RequirePermission(rbac, permissions.JobsUpdate), jobH.UpdateJob)
		jobsGroup.DELETE("/:id", middleware.RequirePermission(rbac, permissions.JobsDelete), jobH.DeleteJob)
		jobsGroup.PATCH("/:id/publish", middleware.RequirePermission(rbac, permissions.JobsPublish), jobH.PublishJob)
		jobsGroup.PATCH("/:id/close", middleware.RequirePermission(rbac, permissions.JobsClose), jobH.CloseJob)
	}

	candidatesGroup := protected.Group("/candidates")
	{
		candidatesGroup.GET("", middleware.RequirePermission(rbac, permissions.CandidatesView), candidateH.GetCandidates)
		candidatesGroup.POST("", middleware.RequirePermission(rbac, permissions.CandidatesCreate), candidateH.CreateCandidate)
		candidatesGroup.GET("/:id", middleware.RequirePermission(rbac, permissions.CandidatesView), candidateH.GetCandidate)
		candidatesGroup.PUT("/:id", middleware.RequirePermission(rbac, permissions.CandidatesUpdate), candidateH.UpdateCandidate)
		candidatesGroup.DELETE("/:id", middleware.RequirePermission(rbac, permissions.CandidatesDelete), candidateH.DeleteCandidate)
		candidatesGroup.POST("/:id/upload-resume", middleware.RequirePermission(rbac, permissions.CandidatesUploadResume), candidateH.UploadResume)
		candidatesGroup.GET("/:id/resume", middleware.RequirePermission(rbac, permissions.CandidatesView), candidateH.DownloadResume)
	}

	applicationsGroup := protected.Group("/applications")
	{
		applicationsGroup.GET("", middleware.RequirePermission(rbac, permissions.ApplicationsView), applicationH.GetApplications)
		applicationsGroup.GET("/by-stage", middleware.RequirePermission(rbac, permissions.ApplicationsView), applicationH.GetApplicationsByStage)
		applicationsGroup.POST("", middleware.RequirePermission(rbac, permissions.ApplicationsCreate), applicationH.CreateApplication)
		applicationsGroup.GET("/:id", middleware.RequirePermission(rbac, permissions.ApplicationsView), applicationH.GetApplication)
		applicationsGroup.PUT("/:id", middleware.RequirePermission(rbac, permissions.ApplicationsUpdate), applicationH.UpdateApplication)
		applicationsGroup.PATCH("/:id/move", middleware.RequirePermission(rbac, permissions.ApplicationsMove), applicationH.MoveApplication)
		applicationsGroup.PATCH("/:id/rate", middleware.RequirePermission(rbac, permissions.ApplicationsRate), applicationH.RateApplication)
		applicationsGroup.DELETE("/:id", middleware.RequirePermission(rbac, permissions.ApplicationsDelete), applicationH.DeleteApplication)
	}
}
