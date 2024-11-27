package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/go-chi/chi/v5"
	"github.com/subvisual/fidl/http/jsend"
	"github.com/subvisual/fidl/types"
	"go.uber.org/zap"
)

// nolint:unused
type envelope map[string]any

func (s *Server) Routes(r chi.Router) {
	r.Get("/fetch/{piece}", s.handleRetrieval)
	r.Get("/banks", s.handleBankList)
}

func (s *Server) handleRetrieval(w http.ResponseWriter, r *http.Request) {
	var params RetrievalParams

	qs := r.URL.Query()
	if err := s.Decode(&params, qs); err != nil {
		s.JSON(w, r, http.StatusBadRequest, err)
		return
	}

	if err := s.Validate.Struct(params); err != nil {
		s.JSON(w, r, http.StatusBadRequest, err)
		return
	}

	verifyEndpoint, err := rebuildBankEndpoint(params.Bank, s.ExternalRoute.BankVerify)
	if err != nil {
		s.JSON(w, r, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	if err := Verify(ctx, verifyEndpoint, s.Wallet, params.Authorization, s.Provider.Cost); err != nil {
		s.JSON(w, r, http.StatusInternalServerError, err)
		return
	}

	piece := chi.URLParam(r, "piece")
	accumulator, r, cleanup := s.Forwarder.tracker.Start(r)
	defer cleanup()

	s.Forwarder.Forward(piece, w, r)

	fil := new(types.FIL)
	fil.Int = new(big.Int).Div(
		new(big.Int).Mul(
			big.NewInt(atomic.LoadInt64(accumulator)),
			s.Provider.Cost.Int,
		),
		big.NewInt(s.Provider.SectorSize),
	)

	redeemEndpoint, err := rebuildBankEndpoint(params.Bank, s.ExternalRoute.BankRedeem)
	if err != nil {
		s.JSON(w, r, http.StatusBadRequest, err)
		return
	}
	if err := Redeem(ctx, redeemEndpoint, s.Wallet, params.Authorization, *fil); err != nil {
		zap.L().Error(
			"failed to reedeem",
			zap.String("authorization", params.Authorization.String()),
			zap.Any("amount", fil), zap.Error(err),
		)
	}

	zap.L().Debug(
		"finished retrieval",
		zap.String("authorization", params.Authorization.String()),
		zap.Any("amount", fil),
		zap.String("piece", piece),
	)
}

func (s *Server) handleBankList(w http.ResponseWriter, r *http.Request) {
	payload := make([]BankListResponse, 0, len(s.Bank))

	for _, v := range s.Bank {
		payload = append(payload, BankListResponse{URL: v.Register, Cost: s.Provider.Cost})
	}

	s.JSON(w, r, http.StatusOK, payload)
}

func (s *Server) HandleForwarderError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, http.ErrHandlerTimeout):
		s.JSON(w, r, http.StatusGatewayTimeout, err)
	default:
		s.JSON(w, r, http.StatusBadGateway, err)
	}
}

func (s *Server) HandleForwarderResponse(resp *http.Response) error {
	trackerID, ok := resp.Request.Context().Value(trackerCtxKey).(string)
	if !ok {
		return nil
	}

	if resp.StatusCode == http.StatusOK {
		resp.Body = &TrackingReader{
			ReadCloser: resp.Body,
			tracker:    s.Forwarder.tracker,
			trackerID:  trackerID,
		}

		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	payload, err := json.Marshal(jsend.Fail(string(body)))
	if err != nil {
		return fmt.Errorf("failed to parse marshal payload: %w", err)
	}

	resp.Body = io.NopCloser(strings.NewReader(string(payload)))
	resp.Header.Set("Content-Type", "application/json")
	resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload)))

	return nil
}
