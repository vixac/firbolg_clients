package local_bullet

import (
	"github.com/vixac/bullet/store/store_interface"
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

func ledgerRecord(record store_interface.LedgerRecord) bullet_interface.LedgerRecord {
	return bullet_interface.LedgerRecord{
		LedgerID:  string(record.LedgerID),
		Position:  int64(record.Position),
		AppendID:  string(record.AppendID),
		CreatedAt: record.CreatedAt,
		Payload:   record.Payload,
	}
}

func ledgerRecords(records []store_interface.LedgerRecord) []bullet_interface.LedgerRecord {
	result := make([]bullet_interface.LedgerRecord, len(records))
	for i, record := range records {
		result[i] = ledgerRecord(record)
	}
	return result
}

func ledgerSelector(selector bullet_interface.LedgerSelector) store_interface.LedgerSelector {
	ids := make([]store_interface.LedgerID, len(selector.LedgerIDs))
	for i, id := range selector.LedgerIDs {
		ids[i] = store_interface.LedgerID(id)
	}
	return store_interface.LedgerSelector{All: selector.All, LedgerIDs: ids, Prefix: selector.Prefix}
}

func (l *LocalBullet) LedgerAppend(req bullet_interface.LedgerAppendRequest) (*bullet_interface.LedgerRecord, error) {
	record, err := l.Store.LedgerAppend(l.Space, store_interface.LedgerID(req.LedgerID), store_interface.LedgerAppendID(req.AppendID), req.Payload)
	if err != nil {
		return nil, err
	}
	result := ledgerRecord(record)
	return &result, nil
}

func (l *LocalBullet) LedgerAppendMany(req bullet_interface.LedgerAppendManyRequest) (*bullet_interface.LedgerAppendManyResponse, error) {
	items := make([]store_interface.LedgerAppendItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = store_interface.LedgerAppendItem{AppendID: store_interface.LedgerAppendID(item.AppendID), Payload: item.Payload}
	}
	records, err := l.Store.LedgerAppendMany(l.Space, store_interface.LedgerID(req.LedgerID), items)
	if err != nil {
		return nil, err
	}
	return &bullet_interface.LedgerAppendManyResponse{Records: ledgerRecords(records)}, nil
}

func (l *LocalBullet) LedgerReadBackward(req bullet_interface.LedgerReadBackwardRequest) (*bullet_interface.LedgerPage, error) {
	page, err := l.Store.LedgerReadBackward(l.Space, ledgerSelector(req.LedgerSelector), req.Cursor, req.Limit)
	if err != nil {
		return nil, err
	}
	return &bullet_interface.LedgerPage{Records: ledgerRecords(page.Records), NextCursor: page.NextCursor}, nil
}

func (l *LocalBullet) LedgerReadForward(req bullet_interface.LedgerReadForwardRequest) (*bullet_interface.LedgerReadForwardResponse, error) {
	var through *store_interface.LedgerPosition
	if req.ThroughPosition != nil {
		value := store_interface.LedgerPosition(*req.ThroughPosition)
		through = &value
	}
	records, err := l.Store.LedgerReadForward(l.Space, ledgerSelector(req.LedgerSelector), store_interface.LedgerPosition(req.AfterPosition), through, req.Limit)
	if err != nil {
		return nil, err
	}
	return &bullet_interface.LedgerReadForwardResponse{Records: ledgerRecords(records)}, nil
}

func (l *LocalBullet) LedgerDelete(req bullet_interface.LedgerDeleteRequest) error {
	return l.Store.LedgerDelete(l.Space, store_interface.LedgerID(req.LedgerID))
}
