package gitlab

type Repository struct {
	Id                int       `json:"id"`
	Name              string    `json:"name"`
	Path              string    `json:"path"`
	PathWithNamespace string    `json:"path_with_namespace"`
	DefaultBranch     string    `json:"default_branch"`
	Namespace         Namespace `json:"namespace"`
	Archived          bool      `json:"archived"`
	HttpUrlToRepo     string    `json:"http_url_to_repo"`
	WebUrl            string    `json:"web_url"`
	CreateAt          string    `json:"created_at"`
	LastActivityAt    string    `json:"last_activity_at"`
	Visibility        string    `json:"visibility"`
	Description       string    `json:"description"`
	CreatorId         int       `json:"creator_id"`
}

// Namespace gitlab namespace
type Namespace struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Kind     string `json:"kind"`
	FullPath string `json:"full_path"` // 全路径
	ParentId int    `json:"parent_id"`
}

// GitlabRepository gitlab repository
type GitlabRepository struct {
	// gorm.Model
	Id                int    `gorm:"cloumn:id"`
	Name              string `gorm:"cloumn:name"`                // name
	Path              string `gorm:"cloumn:path"`                // path
	PathWithNamespace string `gorm:"cloumn:path_with_namespace"` // path with namespace
	DefaultBranch     string `gorm:"cloumn:default_branch"`      // default branch
	Language          string `gorm:"cloumn:language"`            // programming language
	Archived          bool   `gorm:"cloumn:archived"`            // archived
	HttpUrlToRepo     string `gorm:"cloumn:http_url_to_repo"`    // http url to repo
	WebUrl            string `gorm:"cloumn:web_url"`             // web url
	LastActivityAt    string `gorm:"cloumn:last_activity_at"`    // last activity at
	CreateAt          string `gorm:"cloumn:created_at"`          // created at
	Visibility        string `gorm:"cloumn:visibility"`          // visibility
	Description       string `gorm:"cloumn:description"`         // description
	CreatorId         int    `gorm:"cloumn:creator_id"`          // creator id
	NamespaceId       int    `gorm:"cloumn:namespace_id"`        // namespace id
	NamespaceName     string `gorm:"cloumn:namespace_name"`      // namespace name
	NamespacePath     string `gorm:"cloumn:namespace_path"`      // namespace path
	NamespaceKind     string `gorm:"cloumn:namespace_kind"`      // namespace kind
	NamespaceFullPath string `gorm:"cloumn:namespace_full_path"` // namespace full path
	NamespaceParentId int    `gorm:"cloumn:namespace_parent_id"` // namespace parent id
}
