package legacy

import (
	"fmt"
	"net/http"

	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service/httpjson"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type BridgeService struct {
	bridge_service.BridgeServiceBase
	httpClient *http.Client
}

// NewBridgeService creates an instance of the BridgeService.
func NewBridgeService(url string, insecure bool) (*BridgeService, error) {
	return &BridgeService{
		httpClient:        httpjson.NewHTTPClient(insecure),
		BridgeServiceBase: bridge_service.NewBridgeServiceBase(url),
	}, nil
}

func (s *BridgeService) GetDeposit(depositNetwork, depositCount uint32) (*bridge_service.Deposit, error) {
	endpoint := fmt.Sprintf("%s/bridge?net_id=%d&deposit_cnt=%d", s.BridgeServiceBase.Url(), depositNetwork, depositCount)
	resp, _, err := httpjson.HTTPGet[GetDepositResponse](s.httpClient, endpoint)
	if err != nil {
		return nil, err
	}

	if resp.Code != nil {
		errMsg := "unable to retrieve bridge deposit"
		l := log.Warn().Int("code", *resp.Code)
		if resp.Message != nil {
			l.Str("message", *resp.Message)
		}
		l.Msg(errMsg)
		return nil, bridge_service.ErrNotFound
	}

	deposit, err := resp.Deposit.ToDeposit()
	if err != nil {
		return nil, err
	}
	return deposit, nil
}

func (s *BridgeService) GetDeposits(destinationAddress string, offset, limit int) ([]bridge_service.Deposit, int, error) {
	url := fmt.Sprintf("%s/bridges/%s?offset=%d&limit=%d", s.BridgeServiceBase.Url(), destinationAddress, offset, limit)
	resp, _, err := httpjson.HTTPGet[GetDepositsResponse](s.httpClient, url)
	if err != nil {
		return nil, 0, err
	}
	deposits := make([]bridge_service.Deposit, 0, len(resp.Deposits))
	for _, d := range resp.Deposits {
		deposit, err := d.ToDeposit()
		if err != nil {
			return nil, 0, err
		}
		deposits = append(deposits, *deposit)
	}

	return deposits, resp.Total, nil

}

func (s *BridgeService) GetProof(depositNetwork, depositCount uint32, ger *common.Hash) (*bridge_service.Proof, error) {
	endpoint := fmt.Sprintf("%s/merkle-proof?net_id=%d&deposit_cnt=%d", s.BridgeServiceBase.Url(), depositNetwork, depositCount)
	if ger != nil {
		endpoint = fmt.Sprintf("%s/merkle-proof-by-ger?net_id=%d&deposit_cnt=%d&ger=%s", s.BridgeServiceBase.Url(), depositNetwork, depositCount, ger.String())
	}

	resp, _, err := httpjson.HTTPGet[GetProofResponse](s.httpClient, endpoint)
	if err != nil {
		return nil, err
	}

	proof := resp.Proof.ToProof()
	return proof, nil
}
