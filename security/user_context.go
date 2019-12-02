package security

var (
	defaultUserContext = &UserContext{
		UserName: "anonymous",
		Roles:    nil,
		TenantId: "",
		Scopes:   nil,
		Token:    "",
	}
)

type UserContext struct {
	UserName    string   `json:"user_name"`
	Roles       []string `json:"roles"`
	TenantId    string   `json:"tenant_id"`
	Scopes      []string `json:"scope"`
	Authorities []string `json:"authorities"`
	Token       string   `json:"-"`
}

func (c *UserContext) Clone() *UserContext {
	return &UserContext{
		UserName:    c.UserName,
		Roles:       c.Roles[:],
		TenantId:    c.TenantId,
		Scopes:      c.Scopes[:],
		Authorities: c.Authorities[:],
		Token:       c.Token,
	}
}
