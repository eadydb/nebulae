package git

type Repository struct {
	Id                int       `json:"id"`
	Name              string    `json:"name"`
	Path              string    `json:"path"`
	PathWithNamespace string    `json:"path_with_namespace"`
	DefaultBranch     string    `json:"default_branch"`
	Namespace         Namespace `json:"namespace"`
	Archived          bool      `json:"archived"`
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
