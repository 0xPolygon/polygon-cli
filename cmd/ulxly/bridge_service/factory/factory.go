package bridge_service_factory

import (
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service/aggkit"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service/legacy"
)

func NewBridgeService(url string, insecure, useLegacy bool) (bridge_service.BridgeService, error) {
	if useLegacy {
		return legacy.NewBridgeService(url, insecure)
	}
	return aggkit.NewBridgeService(url, insecure)
}
