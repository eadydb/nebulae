package apollo

import (
	"log/slog"
	"testing"
)

var bootstrap = `
apollo:
  meta: ${APOLLO_META:http://service-apollo-config-server-dev.sre.svc.cluster.local:8080}
  bootstrap:
    enabled: true
    namespaces: application.yml,rocket-payment.properties,rocket-business.properties,workflow.config.properties,netty-redis.properties,kf-event-bus.properties,kf-income.properties,sentinel-app-rules.properties,sentinel-dashboard-business.properties,kf-websocket.properties,charging.public-switch.properties,map-common.properties,third-party-info.properties,rabbit-websocket.properties,payment-refactor.properties,customer-backup-domain.properties,lyyProConfig.properties,cache-redis.properties,middle-stage.properties
  eagerLoad:
    enabled: true

app:
  id: charging-business
`

func TestUnmalshalApollo(t *testing.T) {
	namespaces, appId, err := UnmarshalApolloContent([]byte(bootstrap))
	if err != nil {
		t.Errorf("UnmarshalApolloContent error: %v", err)
	}
	if appId != "charging-business" {
		t.Errorf("appId not equal")
	}
	slog.Info("namespaces", slog.Any("namespaces", namespaces))
}
