package committer

import (
	"context"
	"errors"
	"testing"

	"github.com/dawitel/catalog-proj/internal/commitplan"
	commitplanmocks "github.com/dawitel/catalog-proj/mocks/commitplan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCommitter_Apply_Success(t *testing.T) {
	exec := commitplanmocks.NewMockExecutor(t)
	exec.EXPECT().Execute(context.Background(), mock.Anything).Return(nil)
	c := NewWithExecutor(exec)
	plan := commitplan.NewPlan()
	plan.Add(nil)
	err := c.Apply(context.Background(), plan)
	require.NoError(t, err)
}

func TestCommitter_Apply_ExecutorError(t *testing.T) {
	wantErr := errors.New("tx failed")
	exec := commitplanmocks.NewMockExecutor(t)
	exec.EXPECT().Execute(context.Background(), mock.Anything).Return(wantErr)
	c := NewWithExecutor(exec)
	plan := commitplan.NewPlan()
	err := c.Apply(context.Background(), plan)
	assert.ErrorIs(t, err, wantErr)
}
