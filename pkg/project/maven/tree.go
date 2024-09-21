package maven

import (
	"bufio"
	"errors"
	"log/slog"
	"os"
)

const (
	MvnTabSize = 3
	MvnMain    = "+- "
	MvnSub     = "|  "
	MvnLast    = "\\- "
	Suffix     = "deps.txt"
)

type MavenTxt struct {
	PTxt  string
	Txt   string
	Depth int
}

type DependencyTxt struct {
	Project string
	Txts    []MavenTxt
	Path    string
}

// parseMavenDependencyTxt 解析mvn dependency:tree 命令生成的txt文件
func parseMavenDependencyTxt(fileName string) (*DependencyTxt, error) {
	slog.Info("start analyize maven dependency tree", slog.String("fileName", fileName))
	if fileName == "" {
		slog.Error("maven dependency tree file name is empty")
		return nil, errors.New("maven dependency tree file name is empty")
	}
	readFile, err := os.Open(fileName)
	if err != nil {
		slog.Error("read maven dependency tree file failed", slog.String("err", err.Error()), slog.String("fileName", fileName))
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	maven := make([]MavenTxt, 0)
	index := 0
	project := ""
	parents := make([]string, 0)
	preDepth := 0
	for fileScanner.Scan() {
		if index == 0 {
			project = fileScanner.Text()
		} else {
			depth, txt := CalculateDependencyDepth(fileScanner.Text())
			if len(txt) > 0 && txt != "" {

				// 如果深度为0，说明是根节点，需要清空父节点数组
				if depth == 0 {
					parents = parents[:0]
				}

				// 如果深度小于等于上一个深度，说明是上一个节点的父节点，需要删除到上一个节点的父节点
				if depth > 0 && depth <= preDepth {
					parents = parents[:depth]
				}

				mvn := MavenTxt{
					Txt:   txt,
					Depth: depth,
				}
				// 如果父节点数组长度大于0，说明有父节点，设置父节点
				if len(parents) > 0 && depth > 0 {
					mvn.PTxt = parents[depth-1]
				}
				maven = append(maven, mvn)

				parents = append(parents, txt)
				preDepth = depth
			}
		}
		index++
	}

	path := fileName[:len(fileName)-len(Suffix)]
	return &DependencyTxt{
		Project: project,
		Txts:    maven,
		Path:    path,
	}, nil
}

// CalculateDependencyDepth 计算依赖深度
func CalculateDependencyDepth(text string) (int, string) {
	for i := 0; i < len(text)/4; i++ {
		chart := text[i*MvnTabSize : (i+1)*MvnTabSize]
		if chart == MvnMain || chart == MvnLast {
			return i, text[(i+1)*MvnTabSize:]
		}
	}
	return 0, ""
}
