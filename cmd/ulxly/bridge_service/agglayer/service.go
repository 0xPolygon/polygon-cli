package agglayer

import (
	"fmt"
	"net/http"

	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service/httpjson"
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
	endpoint := fmt.Sprintf("%s/bridges?network_id=%d&deposit_count=%d", s.BridgeServiceBase.Url(), depositNetwork, depositCount)
	resp, respError, statusCode, err := httpjson.HTTPGetWithError[GetBridgeResponse, ErrorResponse](s.httpClient, endpoint)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		errMsg := "unable to retrieve bridge deposit"
		log.Warn().Int("code", statusCode).Str("message", respError.Error).Msgf("%s", errMsg)
		return nil, bridge_service.ErrUnableToRetrieveDeposit
	}

	deposit := resp.Bridges[0].ToDeposit()

	return deposit, nil
}

func (s *BridgeService) GetDeposits(destinationAddress string, offset, limit int) ([]bridge_service.Deposit, int, error) {
	pageSize := limit
	pageNumber := offset/limit + 1
	skipItems := offset % limit

	const endpointTemplate = "%s/bridges?from_address=%s&page_number=%d&page_size=%d"

	// loads all deposits when offset is exactly the size of a page or the first part of them when offset is not
	// exactly the size of a page
	endpoint := fmt.Sprintf(endpointTemplate, s.BridgeServiceBase.Url(), destinationAddress, pageNumber, pageSize)
	resp, respError, statusCode, err := httpjson.HTTPGetWithError[GetBridgeResponse, ErrorResponse](s.httpClient, endpoint)
	if err != nil {
		return nil, 0, err
	}

	if statusCode != http.StatusOK {
		errMsg := "unable to retrieve bridge deposits"
		log.Warn().Int("code", statusCode).Str("message", respError.Error).Msgf("%s", errMsg)
		return nil, 0, bridge_service.ErrUnableToRetrieveDeposit
	}

	bridgesResponses := make([]BridgeResponse, 0, limit)
	bridgesResponses = append(bridgesResponses, resp.Bridges[skipItems:pageSize]...)

	// loads the remaining part of deposits when offset is not exactly the size of a page
	// this is needed because the API only supports pagination by page number and page size
	// and not by offset and limit
	if skipItems > 0 {
		endpoint := fmt.Sprintf(endpointTemplate, s.BridgeServiceBase.Url(), destinationAddress, pageNumber+1, pageSize)
		resp, respError, statusCode, err = httpjson.HTTPGetWithError[GetBridgeResponse, ErrorResponse](s.httpClient, endpoint)
		if err != nil {
			return nil, 0, err
		}

		if statusCode != http.StatusOK {
			errMsg := "unable to retrieve bridge deposits"
			log.Warn().Int("code", statusCode).Str("message", respError.Error).Msgf("%s", errMsg)
			return nil, 0, bridge_service.ErrUnableToRetrieveDeposit
		}

		end := skipItems
		if end > len(resp.Bridges) {
			end = len(resp.Bridges)
		}

		bridgesResponses = append(bridgesResponses, resp.Bridges[0:end]...)
	}

	deposits := make([]bridge_service.Deposit, 0, len(bridgesResponses))
	for _, d := range bridgesResponses {
		deposit := d.ToDeposit()
		deposits = append(deposits, *deposit)
	}

	return deposits, resp.Count, nil

}

func (s *BridgeService) GetProof(depositNetwork, depositCount uint32) (*bridge_service.Proof, error) {
	l1InfoTreeIndexEndpoint := fmt.Sprintf("%s/l1-info-tree-index?network_id=%d&deposit_count=%d", s.BridgeServiceBase.Url(), depositNetwork, depositCount)
	l1InfoTreeIndex, l1InfoTreeIndexRespError, statusCode, err := httpjson.HTTPGetWithError[int, ErrorResponse](s.httpClient, l1InfoTreeIndexEndpoint)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		errMsg := "unable to retrieve l1 info tree index"
		log.Warn().Int("code", statusCode).Str("message", l1InfoTreeIndexRespError.Error).Msgf("%s", errMsg)
		return nil, fmt.Errorf(l1InfoTreeIndexRespError.Error)
	}

	endpoint := fmt.Sprintf("%s/claim-proof?network_id=%d&leaf_index=%d&deposit_count=%d", s.BridgeServiceBase.Url(), depositNetwork, l1InfoTreeIndex, depositCount)
	resp, respError, statusCode, err := httpjson.HTTPGetWithError[GetClaimProofResponse, ErrorResponse](s.httpClient, endpoint)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		errMsg := "unable to retrieve bridge deposits"
		log.Warn().Int("code", statusCode).Str("message", respError.Error).Msgf("%s", errMsg)
		return nil, fmt.Errorf(respError.Error)
	}

	proof := resp.ToProof()
	return proof, nil
}
