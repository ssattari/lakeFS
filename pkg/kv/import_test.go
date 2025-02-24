package kv_test

import (
	"bytes"
	"context"
	_ "embed"
	"testing"

	kvpg "github.com/treeverse/lakefs/pkg/kv/postgres"
	"github.com/treeverse/lakefs/pkg/testutil"

	"github.com/stretchr/testify/require"
	"github.com/treeverse/lakefs/pkg/kv"
	_ "github.com/treeverse/lakefs/pkg/kv/mem"
)

func TestImport(t *testing.T) {
	ctx := context.Background()
	kvStore := GetStore(ctx, t)
	defer kvStore.Close()
	tests := []struct {
		name    string
		data    kvpg.MigrateFunc
		err     error
		entries int
	}{
		{
			name:    "basic",
			data:    testutil.MigrateBasic,
			err:     nil,
			entries: 5,
		},
		{
			name: "empty",
			data: testutil.MigrateEmpty,
			err:  kv.ErrInvalidFormat,
		},
		{
			name: "no_header",
			data: testutil.MigrateNoHeader,
			err:  kv.ErrInvalidFormat,
		},
		{
			name: "bad_entry",
			data: testutil.MigrateBadEntry,
			err:  kv.ErrInvalidFormat,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			testutil.MustDo(t, "fill buffer", tt.data(ctx, nil, &buf))
			err := kv.Import(ctx, &buf, kvStore)
			require.ErrorIs(t, err, tt.err)

			testutil.ValidateKV(ctx, t, kvStore, tt.entries)
			testutil.CleanupKV(ctx, t, kvStore)
		})
	}
}
