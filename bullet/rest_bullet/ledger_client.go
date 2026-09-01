package rest_bullet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"github.com/vixac/firbolg_clients/bullet/util"
)

type ledgerRecordWire struct {
	LedgerID  string    `json:"ledger_id"`
	Position  string    `json:"position"`
	AppendID  string    `json:"append_id"`
	CreatedAt time.Time `json:"created_at"`
	Payload   string    `json:"payload"`
}

func convertLedgerRecord(record ledgerRecordWire) (bullet_interface.LedgerRecord, error) {
	position, err := strconv.ParseInt(record.Position, 10, 64)
	if err != nil {
		return bullet_interface.LedgerRecord{}, fmt.Errorf("invalid ledger position %q: %w", record.Position, err)
	}
	return bullet_interface.LedgerRecord{LedgerID: record.LedgerID, Position: position, AppendID: record.AppendID, CreatedAt: record.CreatedAt, Payload: record.Payload}, nil
}

func convertLedgerRecords(records []ledgerRecordWire) ([]bullet_interface.LedgerRecord, error) {
	result := make([]bullet_interface.LedgerRecord, len(records))
	for i, record := range records {
		converted, err := convertLedgerRecord(record)
		if err != nil {
			return nil, err
		}
		result[i] = converted
	}
	return result, nil
}

func decodeLedgerResponse(data []byte, target any) error {
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal ledger response: %w, message body was %q", err, string(data))
	}
	return nil
}

func (c *RestClient) LedgerAppend(req bullet_interface.LedgerAppendRequest) (*bullet_interface.LedgerRecord, error) {
	body, err := util.MarshalJSONBody(req)
	if err != nil {
		return nil, err
	}
	data, err := c.PostReq("/ledger/"+url.PathEscape(req.LedgerID)+"/entries", body, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	var wire ledgerRecordWire
	if err := decodeLedgerResponse(data, &wire); err != nil {
		return nil, err
	}
	result, err := convertLedgerRecord(wire)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *RestClient) LedgerAppendMany(req bullet_interface.LedgerAppendManyRequest) (*bullet_interface.LedgerAppendManyResponse, error) {
	body, err := util.MarshalJSONBody(struct {
		Items []bullet_interface.LedgerAppendItem `json:"items"`
	}{Items: req.Items})
	if err != nil {
		return nil, err
	}
	data, err := c.PostReq("/ledger/"+url.PathEscape(req.LedgerID)+"/entries/batch", body, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	var wire struct {
		Records []ledgerRecordWire `json:"records"`
	}
	if err := decodeLedgerResponse(data, &wire); err != nil {
		return nil, err
	}
	records, err := convertLedgerRecords(wire.Records)
	if err != nil {
		return nil, err
	}
	return &bullet_interface.LedgerAppendManyResponse{Records: records}, nil
}

func (c *RestClient) LedgerReadBackward(req bullet_interface.LedgerReadBackwardRequest) (*bullet_interface.LedgerPage, error) {
	body, err := util.MarshalJSONBody(req)
	if err != nil {
		return nil, err
	}
	data, err := c.PostReq("/ledger/read/backward", body, http.StatusOK)
	if err != nil {
		return nil, err
	}
	var wire struct {
		Records    []ledgerRecordWire `json:"records"`
		NextCursor *string            `json:"next_cursor,omitempty"`
	}
	if err := decodeLedgerResponse(data, &wire); err != nil {
		return nil, err
	}
	records, err := convertLedgerRecords(wire.Records)
	if err != nil {
		return nil, err
	}
	return &bullet_interface.LedgerPage{Records: records, NextCursor: wire.NextCursor}, nil
}

func (c *RestClient) LedgerReadForward(req bullet_interface.LedgerReadForwardRequest) (*bullet_interface.LedgerReadForwardResponse, error) {
	type wireRequest struct {
		bullet_interface.LedgerSelector
		AfterPosition   string  `json:"after_position,omitempty"`
		ThroughPosition *string `json:"through_position,omitempty"`
		Limit           int     `json:"limit"`
	}
	var through *string
	if req.ThroughPosition != nil {
		value := strconv.FormatInt(*req.ThroughPosition, 10)
		through = &value
	}
	body, err := util.MarshalJSONBody(wireRequest{LedgerSelector: req.LedgerSelector, AfterPosition: strconv.FormatInt(req.AfterPosition, 10), ThroughPosition: through, Limit: req.Limit})
	if err != nil {
		return nil, err
	}
	data, err := c.PostReq("/ledger/read/forward", body, http.StatusOK)
	if err != nil {
		return nil, err
	}
	var wire struct {
		Records []ledgerRecordWire `json:"records"`
	}
	if err := decodeLedgerResponse(data, &wire); err != nil {
		return nil, err
	}
	records, err := convertLedgerRecords(wire.Records)
	if err != nil {
		return nil, err
	}
	return &bullet_interface.LedgerReadForwardResponse{Records: records}, nil
}

func (c *RestClient) LedgerDelete(req bullet_interface.LedgerDeleteRequest) error {
	_, err := c.DeleteReq("/ledger/"+url.PathEscape(req.LedgerID), nil, http.StatusNoContent)
	return err
}
