package bullet_interface

import "time"

type LedgerAppendItem struct {
	AppendID string `json:"append_id"`
	Payload  string `json:"payload"`
}

type LedgerAppendRequest struct {
	LedgerID string `json:"-"`
	AppendID string `json:"append_id"`
	Payload  string `json:"payload"`
}

type LedgerAppendManyRequest struct {
	LedgerID string             `json:"-"`
	Items    []LedgerAppendItem `json:"items"`
}

type LedgerSelector struct {
	All       bool     `json:"all"`
	LedgerIDs []string `json:"ledger_ids,omitempty"`
}

type LedgerReadBackwardRequest struct {
	LedgerSelector
	Cursor *string `json:"cursor,omitempty"`
	Limit  int     `json:"limit"`
}

type LedgerReadForwardRequest struct {
	LedgerSelector
	AfterPosition   int64  `json:"-"`
	ThroughPosition *int64 `json:"-"`
	Limit           int    `json:"limit"`
}

type LedgerDeleteRequest struct {
	LedgerID string `json:"-"`
}

type LedgerRecord struct {
	LedgerID  string    `json:"ledger_id"`
	Position  int64     `json:"position"`
	AppendID  string    `json:"append_id"`
	CreatedAt time.Time `json:"created_at"`
	Payload   string    `json:"payload"`
}

type LedgerAppendManyResponse struct {
	Records []LedgerRecord `json:"records"`
}

type LedgerPage struct {
	Records    []LedgerRecord `json:"records"`
	NextCursor *string        `json:"next_cursor,omitempty"`
}

type LedgerReadForwardResponse struct {
	Records []LedgerRecord `json:"records"`
}
