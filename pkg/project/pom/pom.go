package pom

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"os"

	"github.com/eadydb/nebulae/pkg/utils/walk"
)

type Depndency struct {
	GroupId    string
	ArtifactId string
	Version    string
	Scope      string
}

type POM struct {
	XMLName     xml.Name `xml:"project"`
	GroupID     string   `xml:"groupId"`
	ArtifactID  string   `xml:"artifactId"`
	Version     string   `xml:"version"`
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	Parent      struct {
		GroupID    string `xml:"groupId"`
		ArtifactID string `xml:"artifactId"`
		Version    string `xml:"version"`
	} `xml:"parent"`
	Modules []string `xml:"modules>module"`
}

// ParsePOM parses a pom.xml file and returns a POM struct
func ParsePOM(filePath string) (*POM, error) {
	// Read the pom.xml file
	content, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error("failed to read pom.xml", slog.String("error", err.Error()), slog.String("file", filePath))
		return nil, fmt.Errorf("failed to read pom.xml: %w", err)
	}
	return ParsePOMContent(content, filePath)
}

// ParsePOM parses a pom.xml file and returns a POM struct
func ParsePOMContent(content []byte, filePath string) (*POM, error) {
	// Create a POM struct to unmarshal the XML into
	var pom POM

	// Unmarshal the XML into the POM struct
	err := xml.Unmarshal(content, &pom)
	if err != nil {
		slog.Error("failed to parse pom.xml", slog.String("error", err.Error()), slog.String("file", filePath))
		return nil, fmt.Errorf("failed to parse pom.xml: %w", err)
	}

	return &pom, nil
}

// LoadingPomFile 获取pom.xml文件
func LoadingPomFile(path string) ([]string, error) {
	txtPaths, err := walk.From(path).CollectFilterPaths("pom.xml")
	if err != nil {
		slog.Error("loading mvn txt file failed", slog.String("dir", path), slog.String("err", err.Error()), slog.String("fileName", "pom.xml"))
		return nil, err
	}
	return txtPaths, nil
}
