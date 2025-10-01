package aggkit

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service/httpjson"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

const urlPath = "bridge/v1"

type BridgeService struct {
	bridge_service.BridgeServiceBase
	httpClient *http.Client
}

// NewBridgeService creates an instance of the BridgeService.
func NewBridgeService(url string, insecure bool) (*BridgeService, error) {
	return &BridgeService{
		BridgeServiceBase: bridge_service.NewBridgeServiceBase(url),
		httpClient:        httpjson.NewHTTPClient(insecure),
	}, nil
}

func (s *BridgeService) GetDeposit(depositNetwork, depositCount uint32) (*bridge_service.Deposit, error) {
	bridgeEndpoint := fmt.Sprintf("%s/%s/bridges?network_id=%d&deposit_count=%d", s.BridgeServiceBase.Url(), urlPath, depositNetwork, depositCount)
	bridgeResp, bridgeRespError, statusCode, err := httpjson.HTTPGetWithError[getBridgesResponse, errorResponse](s.httpClient, bridgeEndpoint)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		errMsg := "unable to retrieve bridge deposit"
		log.Warn().Int("code", statusCode).Str("message", bridgeRespError.Error).Msg(errMsg)
		return nil, bridge_service.ErrNotFound
	}

	if len(bridgeResp.Bridges) == 0 {
		return nil, bridge_service.ErrNotFound
	}

	deposit := bridgeResp.Bridges[0].ToDeposit(depositNetwork)

	return deposit, nil
}

func (s *BridgeService) GetDeposits(destinationAddress string, offset, limit int) ([]bridge_service.Deposit, int, error) {
	return nil, 0, fmt.Errorf("GetDeposits is not supported by aggkit bridge service yet")
}

func (s *BridgeService) GetProof(depositNetwork, depositCount uint32, ger *common.Hash) (*bridge_service.Proof, error) {
	var l1InfoTreeIndex uint32

	if ger != nil {
		return nil, errors.New("getting proof by ger is not supported yet by Aggkit bridge service")
	}

	timeout := time.After(time.Minute)
out:
	for {
		idx, err := s.getL1InfoTreeIndex(depositNetwork, depositCount)
		if err != nil && !errors.Is(err, bridge_service.ErrNotFound) {
			return nil, err
		} else if err == nil {
			l1InfoTreeIndex = *idx
			break out
		}
		select {
		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for l1 info tree index")
		default:
			time.Sleep(time.Second)
		}
	}

	endpoint := fmt.Sprintf("%s/%s/claim-proof?network_id=%d&leaf_index=%d&deposit_count=%d", s.BridgeServiceBase.Url(), urlPath, depositNetwork, l1InfoTreeIndex, depositCount)
	resp, respError, statusCode, err := httpjson.HTTPGetWithError[getClaimProofResponse, errorResponse](s.httpClient, endpoint)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		if statusCode == http.StatusNotFound {
			return nil, bridge_service.ErrNotFound
		}
		errMsg := "unable to retrieve proof"
		log.Warn().Int("code", statusCode).Str("message", respError.Error).Msg(errMsg)
		return nil, fmt.Errorf(respError.Error)
	}

	proof := resp.ToProof()
	return proof, nil
}

func (s *BridgeService) getL1InfoTreeIndex(depositNetwork, depositCount uint32) (*uint32, error) {
	l1InfoTreeIndexEndpoint := fmt.Sprintf("%s/%s/l1-info-tree-index?network_id=%d&deposit_count=%d", s.BridgeServiceBase.Url(), urlPath, depositNetwork, depositCount)
	l1InfoTreeIndex, l1InfoTreeIndexRespError, statusCode, err := httpjson.HTTPGetWithError[uint32, errorResponse](s.httpClient, l1InfoTreeIndexEndpoint)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		if statusCode == http.StatusNotFound {
			return nil, bridge_service.ErrNotFound
		}
		if statusCode == http.StatusInternalServerError {
			if strings.HasSuffix(l1InfoTreeIndexRespError.Error, "error: this bridge has not been included on the L1 Info Tree yet") ||
				strings.HasSuffix(l1InfoTreeIndexRespError.Error, "error: not found") {
				return nil, bridge_service.ErrNotFound
			}
		}
		errMsg := "unable to retrieve l1 info tree index"
		log.Warn().Int("code", statusCode).Str("message", l1InfoTreeIndexRespError.Error).Msg(errMsg)
		return nil, fmt.Errorf(l1InfoTreeIndexRespError.Error)
	}

	return &l1InfoTreeIndex, nil
}
