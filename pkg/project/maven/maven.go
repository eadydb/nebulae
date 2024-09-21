package maven

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"time"

	"github.com/eadydb/nebulae/pkg/utils/cmd"
	"github.com/eadydb/nebulae/pkg/walk"
)

type Maven struct {
	Path     string // maven pom.xml path
	Project  string // project name
	FileName string // mvn dependency:tree 生成的txt文件名
	Ctx      context.Context
}

func NewMaven(ctx context.Context, path, project string) *Maven {
	return &Maven{
		Path:    path,
		Project: project,
		Ctx:     ctx,
	}
}

// MavenErrorProject maven build error project
var MavenErrorProject []string = make([]string, 0)

// ExectueCmd 执行mvn dependency:tree 命令
func (m *Maven) ExectueCmd(retry bool) error {
	// mvn dependency:tree -Ppord -DappendOutput=true -DoutputFile=deps.txt -DoutputType=text
	command := exec.CommandContext(m.Ctx, "mvn", "dependency:tree", "-Ppord", "-DappendOutput=true", fmt.Sprintf("-DoutputFile=%s", m.FileName), "-DoutputType=text")
	if body, err := cmd.DefaultExecCommand.RunCmdOutOnce(m.Ctx, command); err != nil {
		if retry {
			// 延迟5秒后重试
			time.Sleep(5 * time.Second)
			slog.Info("retry execute maven dependency:tree cmd", slog.String("project", m.Project), slog.String("path", m.Path))
			return m.ExectueCmd(false)
		} else {
			slog.Error("mvn dependency:tree failed", slog.String("project", m.Project), slog.String("path", m.Path), slog.String("body", string(body)), slog.String("err", err.Error()))
			MavenErrorProject = append(MavenErrorProject, m.Project)
			return err
		}
	}
	return nil
}

// loadingMvnTxtFile 获取mvn dependency:tree 命令生成的txt文件
func (m *Maven) LoadingMvnTxtFile() ([]string, error) {
	txtPaths, err := walk.From(m.Path).CollectFilterPaths(m.FileName)
	if err != nil {
		slog.Error("loading mvn txt file failed", slog.String("dir", m.Path), slog.String("err", err.Error()), slog.String("fileName", m.FileName))
		return nil, err
	}
	return txtPaths, nil
}
