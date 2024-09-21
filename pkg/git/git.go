package git

import "context"

type Gitlab struct {
	PrivateToken string
	Domain       string
	Ctx          context.Context
}

func NewGitlab(ctx context.Context, privateToken, domain string) *Gitlab {
	return &Gitlab{
		PrivateToken: privateToken,
		Domain:       domain,
		Ctx:          ctx,
	}
}

// ScanGitlabHub scan gitlab hub
func (g *Gitlab) ScanGitlabHub() ([]Repository, error) {
	return nil, nil
}

// Clone clone gitlab repository
func (g *Gitlab) Clone(hubDir string) error {
	return nil
}
