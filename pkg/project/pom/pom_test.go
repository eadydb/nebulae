package pom

import (
	"log/slog"
	"testing"
)

var pomXml = `
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
        xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    <groupId>cn.lyy</groupId>
    <artifactId>ic_card_server</artifactId>
    <version>2.0.0-BETA</version>
	<name>ic_card_server</name>
    <packaging>pom</packaging>

    <modules>
        <module>ic_card_service_api</module>
    </modules>

    <parent>
        <groupId>com.lyy.base</groupId>
        <artifactId>base-parent</artifactId>
        <version>1.1.0</version>
    </parent>
</project>
`

func TestPomXml(t *testing.T) {
	pom, err := ParsePOMContent([]byte(pomXml), "pom.xml")
	if err != nil {
		slog.Error("parse pom.xml failed", slog.String("err", err.Error()))
	}
	slog.Info("parse pom.xml success", slog.Any("pom", pom))
}
