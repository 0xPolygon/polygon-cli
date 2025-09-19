package agglayer

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service/httpjson"
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
		httpClient:        httpjson.NewHTTPClient(insecure),
		BridgeServiceBase: bridge_service.NewBridgeServiceBase(url),
	}, nil
}

func (s *BridgeService) GetDeposit(depositNetwork, depositCount uint32) (*bridge_service.Deposit, error) {
	bridgeEndpoint := fmt.Sprintf("%s/%s/bridges?network_id=%d&deposit_count=%d", s.BridgeServiceBase.Url(), urlPath, depositNetwork, depositCount)
	bridgeResp, bridgeRespError, statusCode, err := httpjson.HTTPGetWithError[GetBridgeResponse, ErrorResponse](s.httpClient, bridgeEndpoint)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		errMsg := "unable to retrieve bridge deposit"
		log.Warn().Int("code", statusCode).Str("message", bridgeRespError.Error).Msgf("%s", errMsg)
		return nil, bridge_service.ErrNotFound
	}

	if len(bridgeResp.Bridges) == 0 {
		return nil, bridge_service.ErrNotFound
	}

	deposit, err := s.responseToDeposit(bridgeResp.Bridges[0])
	if err != nil {
		return nil, err
	}

	return deposit, nil
}

func (s *BridgeService) GetDeposits(destinationAddress string, offset, limit int) ([]bridge_service.Deposit, int, error) {
	pageSize := limit
	pageNumber := offset/limit + 1
	skipItems := offset % limit

	const endpointTemplate = "%s/%s/bridges?from_address=%s&page_number=%d&page_size=%d"

	// loads all deposits when offset is exactly the size of a page or the first part of them when offset is not
	// exactly the size of a page
	endpoint := fmt.Sprintf(endpointTemplate, s.BridgeServiceBase.Url(), urlPath, destinationAddress, pageNumber, pageSize)
	resp, respError, statusCode, err := httpjson.HTTPGetWithError[GetBridgeResponse, ErrorResponse](s.httpClient, endpoint)
	if err != nil {
		return nil, 0, err
	}

	if statusCode != http.StatusOK {
		errMsg := "unable to retrieve bridge deposits"
		log.Warn().Int("code", statusCode).Str("message", respError.Error).Msgf("%s", errMsg)
		return nil, 0, bridge_service.ErrNotFound
	}

	bridgesResponses := make([]BridgeResponse, 0, limit)
	bridgesResponses = append(bridgesResponses, resp.Bridges[skipItems:pageSize]...)

	// loads the remaining part of deposits when offset is not exactly the size of a page
	// this is needed because the API only supports pagination by page number and page size
	// and not by offset and limit
	if skipItems > 0 {
		endpoint := fmt.Sprintf(endpointTemplate, s.BridgeServiceBase.Url(), urlPath, destinationAddress, pageNumber+1, pageSize)
		resp, respError, statusCode, err = httpjson.HTTPGetWithError[GetBridgeResponse, ErrorResponse](s.httpClient, endpoint)
		if err != nil {
			return nil, 0, err
		}

		if statusCode != http.StatusOK {
			errMsg := "unable to retrieve bridge deposits"
			log.Warn().Int("code", statusCode).Str("message", respError.Error).Msgf("%s", errMsg)
			return nil, 0, bridge_service.ErrNotFound
		}

		end := skipItems
		if end > len(resp.Bridges) {
			end = len(resp.Bridges)
		}

		bridgesResponses = append(bridgesResponses, resp.Bridges[0:end]...)
	}

	deposits := make([]bridge_service.Deposit, 0, len(bridgesResponses))
	for _, bridgeResp := range bridgesResponses {
		deposit, err := s.responseToDeposit(bridgeResp)
		if err != nil {
			return nil, 0, err
		}
		deposits = append(deposits, *deposit)
	}

	return deposits, resp.Count, nil

}

func (s *BridgeService) GetProof(depositNetwork, depositCount uint32) (*bridge_service.Proof, error) {
	var l1InfoTreeIndex uint32
out:
	for {
		select {
		case <-time.After(time.Minute):
			return nil, fmt.Errorf("timeout waiting for l1 info tree index")
		default:
			idx, err := s.getL1InfoTreeIndex(depositNetwork, depositCount)
			if errors.Is(err, bridge_service.ErrNotFound) {
				time.Sleep(time.Second)
				continue
			}
			if err != nil {
				return nil, err
			}

			l1InfoTreeIndex = *idx
			break out
		}
	}

	endpoint := fmt.Sprintf("%s/%s/claim-proof?network_id=%d&leaf_index=%d&deposit_count=%d", s.BridgeServiceBase.Url(), urlPath, depositNetwork, l1InfoTreeIndex, depositCount)
	resp, respError, statusCode, err := httpjson.HTTPGetWithError[GetClaimProofResponse, ErrorResponse](s.httpClient, endpoint)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		if statusCode == http.StatusNotFound {
			return nil, bridge_service.ErrNotFound
		}
		errMsg := "unable to retrieve proof"
		log.Warn().Int("code", statusCode).Str("message", respError.Error).Msgf("%s", errMsg)
		return nil, fmt.Errorf(respError.Error)
	}

	proof := resp.ToProof()
	return proof, nil
}

func (s *BridgeService) getInjectedL1InfoLeaf(depositNetwork, l1InfoTreeIndex uint32) (*getInjectedL1InfoLeafResponse, error) {
	endpoint := fmt.Sprintf("%s/%s/injected-l1-info-leaf?network_id=%d&leaf_index=%d", s.BridgeServiceBase.Url(), urlPath, depositNetwork, l1InfoTreeIndex)
	resp, errorResp, statusCode, err := httpjson.HTTPGetWithError[*getInjectedL1InfoLeafResponse, ErrorResponse](s.httpClient, endpoint)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		if statusCode == http.StatusNotFound {
			return nil, bridge_service.ErrNotFound
		}
		errMsg := "unable to retrieve l1 info leaf"
		log.Warn().Int("code", statusCode).Str("message", errorResp.Error).Msgf("%s", errMsg)
		return nil, bridge_service.ErrNotFound
	}

	return resp, nil
}

func (s *BridgeService) getL1InfoTreeIndex(depositNetwork, depositCount uint32) (*uint32, error) {
	l1InfoTreeIndexEndpoint := fmt.Sprintf("%s/%s/l1-info-tree-index?network_id=%d&deposit_count=%d", s.BridgeServiceBase.Url(), urlPath, depositNetwork, depositCount)
	l1InfoTreeIndex, l1InfoTreeIndexRespError, statusCode, err := httpjson.HTTPGetWithError[uint32, ErrorResponse](s.httpClient, l1InfoTreeIndexEndpoint)
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
		log.Warn().Int("code", statusCode).Str("message", l1InfoTreeIndexRespError.Error).Msgf("%s", errMsg)
		return nil, fmt.Errorf(l1InfoTreeIndexRespError.Error)
	}

	return &l1InfoTreeIndex, nil
}

func (s *BridgeService) responseToDeposit(bridgeResp BridgeResponse) (*bridge_service.Deposit, error) {
	depositNetwork := bridgeResp.OrigNet
	depositCount := bridgeResp.DepositCnt

	isReadyForClaim := false
	l1InfoTreeIndex, err := s.getL1InfoTreeIndex(depositNetwork, depositCount)
	if err != nil {
		return nil, err
	}
	if l1InfoTreeIndex == nil {
		l1InfoLeaf, iErr := s.getInjectedL1InfoLeaf(depositNetwork, *l1InfoTreeIndex)
		if iErr != nil {
			return nil, iErr
		}
		isReadyForClaim = l1InfoLeaf != nil
	}

	deposit, err := bridgeResp.ToDeposit(isReadyForClaim)
	if err != nil {
		return nil, err
	}

	return deposit, nil
}
