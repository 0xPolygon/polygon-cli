package agglayer

import (
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service"
)

type BridgeService struct {
	bridge_service.BridgeServiceBase
}

// NewBridgeService creates an instance of the BridgeService.
func NewBridgeService(url string, insecure bool) (*BridgeService, error) {
	return &BridgeService{}, nil
}

func (s *BridgeService) GetDeposit(depositNetwork, depositCount uint32) (*bridge_service.Deposit, error) {
	panic("not implemented") // TODO: Implement
}
func (s *BridgeService) GetDeposits(destinationAddress string, offset, limit int) ([]bridge_service.Deposit, int, error) {
	panic("not implemented") // TODO: Implement
}
func (s *BridgeService) GetProof(depositNetwork, depositCount uint32) (*bridge_service.Proof, error) {
	panic("not implemented") // TODO: Implement
}
