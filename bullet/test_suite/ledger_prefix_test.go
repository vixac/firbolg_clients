package test_suite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

func TestLedgerPrefix(t *testing.T) {
	for _, pair := range buildClientPairs(t) {
		t.Run(pair.name, func(t *testing.T) {
			for _, id := range []string{"orders_", "orders_a", "ordersXa", "Orders_a", "orders_b", "other"} {
				_, err := pair.local.LedgerAppend(bullet_interface.LedgerAppendRequest{LedgerID: id, AppendID: "one", Payload: id})
				require.NoError(t, err)
			}
			for name, client := range map[string]bullet_interface.BulletClientInterface{"local": pair.local, "rest": pair.rest} {
				t.Run(name, func(t *testing.T) {
					selector := bullet_interface.LedgerSelector{Prefix: "orders_"}
					page, err := client.LedgerReadBackward(bullet_interface.LedgerReadBackwardRequest{LedgerSelector: selector, Limit: 2})
					require.NoError(t, err)
					require.Len(t, page.Records, 2)
					assert.Equal(t, "orders_b", page.Records[0].LedgerID)
					assert.Equal(t, "orders_a", page.Records[1].LedgerID)
					require.NotNil(t, page.NextCursor)

					next, err := client.LedgerReadBackward(bullet_interface.LedgerReadBackwardRequest{LedgerSelector: selector, Cursor: page.NextCursor, Limit: 2})
					require.NoError(t, err)
					require.Len(t, next.Records, 1)
					assert.Equal(t, "orders_", next.Records[0].LedgerID)
					assert.Nil(t, next.NextCursor)

					_, err = client.LedgerReadBackward(bullet_interface.LedgerReadBackwardRequest{LedgerSelector: bullet_interface.LedgerSelector{Prefix: "orders"}, Cursor: page.NextCursor, Limit: 2})
					require.Error(t, err)

					forward, err := client.LedgerReadForward(bullet_interface.LedgerReadForwardRequest{LedgerSelector: selector, Limit: 10})
					require.NoError(t, err)
					require.Len(t, forward.Records, 3)
					assert.Equal(t, "orders_", forward.Records[0].LedgerID)
					assert.Equal(t, "orders_a", forward.Records[1].LedgerID)
					assert.Equal(t, "orders_b", forward.Records[2].LedgerID)

					through := forward.Records[1].Position
					bounded, err := client.LedgerReadForward(bullet_interface.LedgerReadForwardRequest{LedgerSelector: selector, AfterPosition: forward.Records[0].Position, ThroughPosition: &through, Limit: 10})
					require.NoError(t, err)
					assert.Equal(t, forward.Records[1:2], bounded.Records)
					limited, err := client.LedgerReadForward(bullet_interface.LedgerReadForwardRequest{LedgerSelector: selector, Limit: 1})
					require.NoError(t, err)
					assert.Equal(t, forward.Records[:1], limited.Records)

					missing := bullet_interface.LedgerSelector{Prefix: "missing"}
					emptyForward, err := client.LedgerReadForward(bullet_interface.LedgerReadForwardRequest{LedgerSelector: missing, Limit: 10})
					require.NoError(t, err)
					assert.Empty(t, emptyForward.Records)
					emptyBackward, err := client.LedgerReadBackward(bullet_interface.LedgerReadBackwardRequest{LedgerSelector: missing, Limit: 10})
					require.NoError(t, err)
					assert.Empty(t, emptyBackward.Records)
					assert.Nil(t, emptyBackward.NextCursor)

					for _, invalid := range []bullet_interface.LedgerSelector{
						{Prefix: ""},
						{All: true, Prefix: "orders_"},
						{LedgerIDs: []string{"orders_"}, Prefix: "orders_"},
						{Prefix: "bad%prefix"},
					} {
						_, err := client.LedgerReadBackward(bullet_interface.LedgerReadBackwardRequest{LedgerSelector: invalid, Limit: 10})
						require.Error(t, err, "selector: %+v", invalid)
						_, err = client.LedgerReadForward(bullet_interface.LedgerReadForwardRequest{LedgerSelector: invalid, Limit: 10})
						require.Error(t, err, "selector: %+v", invalid)
					}
				})
			}
		})
	}
}
