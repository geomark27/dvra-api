package identity

const (
	ContextIdentity = "identity"
	ContextTenant   = "tenant"
	ContextPlatform = "platform"
)

// Context identifies the authorization boundary selected for a token.
type Context struct {
	Type         string
	CompanyID    *uint
	MembershipID *uint
	TenantRoleID *uint
	Role         string
}
