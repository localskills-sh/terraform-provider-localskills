package client

// --- Skills ---

type Skill struct {
	ID             string   `json:"id"`
	PublicID       string   `json:"publicId"`
	TenantID       string   `json:"tenantId"`
	Name           string   `json:"name"`
	Slug           string   `json:"slug"`
	Description    string   `json:"description"`
	Type           string   `json:"type"`
	Visibility     string   `json:"visibility"`
	Tags           []string `json:"tags"`
	CurrentVersion int      `json:"currentVersion"`
	CurrentSemver  string   `json:"currentSemver"`
	CreatedBy      string   `json:"createdBy"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}

type CreateSkillRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Type        string   `json:"type"`
	Visibility  string   `json:"visibility"`
	Content     string   `json:"content"`
	TenantID    string   `json:"tenantId"`
	Tags        []string `json:"tags,omitempty"`
}

type UpdateSkillRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Visibility  *string  `json:"visibility,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type SkillVersion struct {
	ID          string `json:"id"`
	SkillID     string `json:"skillId"`
	Version     int    `json:"version"`
	Semver      string `json:"semver"`
	ContentKey  string `json:"contentKey"`
	ContentHash string `json:"contentHash"`
	Message     string `json:"message"`
	Format      string `json:"format"`
	FileCount   int    `json:"fileCount"`
	CreatedBy   string `json:"createdBy"`
	CreatedAt   string `json:"createdAt"`
}

type CreateSkillVersionRequest struct {
	Content string `json:"content"`
	Message string `json:"message,omitempty"`
	Semver  string `json:"semver,omitempty"`
	Bump    string `json:"bump,omitempty"`
}

type SkillContent struct {
	Skill   Skill  `json:"skill"`
	Format  string `json:"format"`
	Content string `json:"content"`
	Version int    `json:"version"`
	Semver  string `json:"semver"`
}

type SkillAnalytics struct {
	TotalDownloads int `json:"totalDownloads"`
	UniqueUsers    int `json:"uniqueUsers"`
	UniqueIPs      int `json:"uniqueIps"`
}

type PackageManifest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Files       []string `json:"files"`
}

type SkillWithVersion struct {
	Skill
	CurrentVersionInfo *SkillVersion `json:"currentVersionInfo"`
}

// --- Tenants (Teams) ---

type Tenant struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type TenantWithRole struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Role        string `json:"role"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type CreateTenantRequest struct {
	Name string `json:"name"`
}

type UpdateTenantRequest struct {
	Name        *string `json:"name,omitempty"`
	Slug        *string `json:"slug,omitempty"`
	Description *string `json:"description,omitempty"`
}

// --- Invitations ---

type TenantInvitation struct {
	ID         string  `json:"id"`
	TenantID   string  `json:"tenantId"`
	Email      string  `json:"email"`
	Role       string  `json:"role"`
	Token      string  `json:"token"`
	InvitedBy  string  `json:"invitedBy"`
	ExpiresAt  string  `json:"expiresAt"`
	AcceptedAt *string `json:"acceptedAt"`
	CreatedAt  string  `json:"createdAt"`
}

type CreateInvitationRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// --- API Tokens ---

type ApiToken struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	LastUsedAt *string `json:"lastUsedAt"`
	ExpiresAt  *string `json:"expiresAt"`
	CreatedAt  string  `json:"createdAt"`
}

type ApiTokenWithSecret struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type TeamApiToken struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	CreatedAt      string  `json:"createdAt"`
	LastUsedAt     *string `json:"lastUsedAt"`
	ExpiresAt      *string `json:"expiresAt"`
	CreatedByName  *string `json:"createdByName"`
	CreatedByEmail string  `json:"createdByEmail"`
}

type TeamApiTokenWithSecret struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Token     string  `json:"token"`
	CreatedAt string  `json:"createdAt"`
	ExpiresAt *string `json:"expiresAt"`
}

type CreateTokenRequest struct {
	Name string `json:"name"`
}

type CreateTeamTokenRequest struct {
	Name          string `json:"name"`
	ExpiresInDays *int   `json:"expiresInDays,omitempty"`
}

// --- SCIM Tokens ---

type ScimToken struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	CreatedAt  string  `json:"createdAt"`
	LastUsedAt *string `json:"lastUsedAt"`
	ExpiresAt  *string `json:"expiresAt"`
}

type ScimTokenWithSecret struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Token     string  `json:"token"`
	CreatedAt string  `json:"createdAt"`
	ExpiresAt *string `json:"expiresAt"`
}

type CreateScimTokenRequest struct {
	Name          string `json:"name"`
	ExpiresInDays *int   `json:"expiresInDays,omitempty"`
}

// --- OIDC Trust Policies ---

type OidcTrustPolicy struct {
	ID                string   `json:"id"`
	TenantID          string   `json:"tenantId"`
	Name              string   `json:"name"`
	Provider          string   `json:"provider"`
	Repository        string   `json:"repository"`
	RefFilter         string   `json:"refFilter"`
	EnvironmentFilter *string  `json:"environmentFilter"`
	SkillIDs          []string `json:"skillIds"`
	Enabled           bool     `json:"enabled"`
	CreatedBy         string   `json:"createdBy"`
	CreatedAt         string   `json:"createdAt"`
	UpdatedAt         string   `json:"updatedAt"`
}

type CreateOidcPolicyRequest struct {
	Name              string   `json:"name"`
	Provider          string   `json:"provider"`
	Repository        string   `json:"repository"`
	RefFilter         string   `json:"refFilter,omitempty"`
	EnvironmentFilter *string  `json:"environmentFilter,omitempty"`
	SkillIDs          []string `json:"skillIds,omitempty"`
	Enabled           bool     `json:"enabled"`
}

type UpdateOidcPolicyRequest struct {
	Name              *string  `json:"name,omitempty"`
	Provider          *string  `json:"provider,omitempty"`
	Repository        *string  `json:"repository,omitempty"`
	RefFilter         *string  `json:"refFilter,omitempty"`
	EnvironmentFilter *string  `json:"environmentFilter,omitempty"`
	SkillIDs          []string `json:"skillIds,omitempty"`
	Enabled           *bool    `json:"enabled,omitempty"`
}

// --- SSO ---

type SsoConnection struct {
	ID           string   `json:"id"`
	TenantID     string   `json:"tenantId"`
	DisplayName  string   `json:"displayName"`
	IdpEntityID  string   `json:"idpEntityId"`
	IdpSsoURL    string   `json:"idpSsoUrl"`
	IdpSloURL    string   `json:"idpSloUrl"`
	SpEntityID   string   `json:"spEntityId"`
	SpAcsURL     string   `json:"spAcsUrl"`
	DefaultRole  string   `json:"defaultRole"`
	EmailDomains []string `json:"emailDomains"`
	Enabled      bool     `json:"enabled"`
	RequireSso   bool     `json:"requireSso"`
	MetadataURL  string   `json:"metadataUrl"`
	CreatedAt    string   `json:"createdAt"`
	UpdatedAt    string   `json:"updatedAt"`
}

type UpdateSsoRequest struct {
	DisplayName  *string  `json:"displayName,omitempty"`
	MetadataURL  *string  `json:"metadataUrl,omitempty"`
	MetadataXML  *string  `json:"metadataXml,omitempty"`
	DefaultRole  *string  `json:"defaultRole,omitempty"`
	EmailDomains []string `json:"emailDomains,omitempty"`
	Enabled      *bool    `json:"enabled,omitempty"`
	RequireSso   *bool    `json:"requireSso,omitempty"`
}

// --- User Profile ---

type UserProfile struct {
	ID       string  `json:"id"`
	Username *string `json:"username"`
	Name     *string `json:"name"`
	Email    string  `json:"email"`
	Image    *string `json:"image"`
	Bio      *string `json:"bio"`
}

// --- Audit Logs ---

type AuditLogEntry struct {
	ID           string  `json:"id"`
	Action       string  `json:"action"`
	ActorID      *string `json:"actorId"`
	ActorName    *string `json:"actorName"`
	ActorImage   *string `json:"actorImage"`
	ResourceType string  `json:"resourceType"`
	ResourceID   string  `json:"resourceId"`
	Metadata     string  `json:"metadata"`
	CreatedAt    string  `json:"createdAt"`
}

type AuditLogResponse struct {
	Entries []AuditLogEntry `json:"entries"`
	Total   int             `json:"total"`
	Page    int             `json:"page"`
}

// --- Explore ---

type ExploreSkill struct {
	ID             string   `json:"id"`
	PublicID       string   `json:"publicId"`
	Name           string   `json:"name"`
	Slug           string   `json:"slug"`
	Description    string   `json:"description"`
	Type           string   `json:"type"`
	Tags           []string `json:"tags"`
	CurrentVersion int      `json:"currentVersion"`
	CurrentSemver  string   `json:"currentSemver"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
	AuthorName     string   `json:"authorName"`
	AuthorUsername string   `json:"authorUsername"`
	Downloads      int      `json:"downloads"`
}
