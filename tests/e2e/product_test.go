package e2e

import (
	"context"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"

	"github.com/dawitel/product-catalog-service/internal/app/product/contracts"
	"github.com/dawitel/product-catalog-service/internal/app/product/domain"
	"github.com/dawitel/product-catalog-service/internal/commitplan"
	"github.com/dawitel/product-catalog-service/internal/app/product/queries/get_product"
	"github.com/dawitel/product-catalog-service/internal/app/product/queries/list_products"
	"github.com/dawitel/product-catalog-service/internal/app/product/repo"
	activate_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/activate_product"
	apply_discount "github.com/dawitel/product-catalog-service/internal/app/product/usecases/apply_discount"
	archive_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/archive_product"
	create_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/create_product"
	deactivate_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/deactivate_product"
	remove_discount "github.com/dawitel/product-catalog-service/internal/app/product/usecases/remove_discount"
	update_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/update_product"
	"github.com/dawitel/product-catalog-service/internal/pkg/config"
	"github.com/dawitel/product-catalog-service/internal/models/m_outbox"
	"github.com/dawitel/product-catalog-service/internal/pkg/clock"
	"github.com/dawitel/product-catalog-service/internal/pkg/committer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (context.Context, *spanner.Client, *create_product.Interactor, *update_product.Interactor, *activate_product.Interactor, *deactivate_product.Interactor, *archive_product.Interactor, *apply_discount.Interactor, *remove_discount.Interactor, *get_product.Query, *list_products.Query) {
	t.Helper()
	cfg := config.LoadFromEnv()
	if cfg.SpannerEmulatorHost == "" {
		t.Skip("SPANNER_EMULATOR_HOST not set, skipping e2e")
	}
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, cfg.DatabasePath())
	require.NoError(t, err)
	t.Cleanup(func() { client.Close() })

	comm := committer.New(client)
	clk := clock.Real{}
	productRepo := repo.NewProductRepo(client)
	outboxRepo := repo.NewOutboxRepo()
	readModel := repo.NewReadModel(client)

	createUC := create_product.New(productRepo, outboxRepo, comm, clk)
	updateUC := update_product.New(productRepo, outboxRepo, comm, clk)
	activateUC := activate_product.New(productRepo, outboxRepo, comm, clk)
	deactivateUC := deactivate_product.New(productRepo, outboxRepo, comm, clk)
	archiveUC := archive_product.New(productRepo, outboxRepo, comm, clk)
	applyDiscUC := apply_discount.New(productRepo, outboxRepo, comm, clk)
	removeDiscUC := remove_discount.New(productRepo, outboxRepo, comm, clk)
	getQuery := get_product.New(readModel, clk)
	listQuery := list_products.New(readModel, clk)

	return ctx, client, createUC, updateUC, activateUC, deactivateUC, archiveUC, applyDiscUC, removeDiscUC, getQuery, listQuery
}

func getOutboxEvents(t *testing.T, client *spanner.Client, aggregateID string) []m_outbox.Row {
	t.Helper()
	ctx := context.Background()
	stmt := spanner.Statement{
		SQL:    "SELECT event_id, event_type, aggregate_id, TO_JSON_STRING(payload) AS payload, status, created_at, processed_at FROM outbox_events WHERE aggregate_id = @aid ORDER BY created_at",
		Params: map[string]interface{}{"aid": aggregateID},
	}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	var rows []m_outbox.Row
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		require.NoError(t, err)
		var e m_outbox.Row
		var processedAt spanner.NullTime
		err = row.Columns(&e.EventID, &e.EventType, &e.AggregateID, &e.Payload, &e.Status, &e.CreatedAt, &processedAt)
		require.NoError(t, err)
		if processedAt.Valid {
			e.ProcessedAt = processedAt.Time
		}
		rows = append(rows, e)
	}
	return rows
}

func TestProductCreationFlow(t *testing.T) {
	ctx, client, createUC, _, _, _, _, _, _, getQuery, _ := setup(t)
	productID, err := createUC.Execute(ctx, create_product.Request{
		Name:                 "Test Product",
		Description:          "Desc",
		Category:             "cat1",
		BasePriceNumerator:   1999,
		BasePriceDenominator: 100,
	})
	require.NoError(t, err)
	require.NotEmpty(t, productID)

	product, err := getQuery.Execute(ctx, productID)
	require.NoError(t, err)
	require.NotNil(t, product)
	assert.Equal(t, "Test Product", product.Name)
	assert.Equal(t, int64(1999), product.EffectivePriceNumerator)
	assert.Equal(t, int64(100), product.EffectivePriceDenominator)

	events := getOutboxEvents(t, client, productID)
	require.Len(t, events, 1)
	assert.Equal(t, "product.created", events[0].EventType)
}

func TestProductUpdateFlow(t *testing.T) {
	ctx, _, createUC, updateUC, _, _, _, _, _, getQuery, _ := setup(t)
	productID, err := createUC.Execute(ctx, create_product.Request{
		Name:                 "Original",
		Description:          "D",
		Category:             "c1",
		BasePriceNumerator:   1000,
		BasePriceDenominator: 100,
	})
	require.NoError(t, err)

	err = updateUC.Execute(ctx, update_product.Request{
		ProductID:   productID,
		Name:        "Updated Name",
		Description: "New desc",
		Category:    "c2",
	})
	require.NoError(t, err)

	product, err := getQuery.Execute(ctx, productID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", product.Name)
	assert.Equal(t, "New desc", product.Description)
	assert.Equal(t, "c2", product.Category)
}

func TestDiscountApplicationFlow(t *testing.T) {
	ctx, _, createUC, _, _, _, _, applyDiscUC, _, getQuery, _ := setup(t)
	productID, err := createUC.Execute(ctx, create_product.Request{
		Name:                 "With Discount",
		Description:          "",
		Category:             "c1",
		BasePriceNumerator:   10000,
		BasePriceDenominator: 100,
	})
	require.NoError(t, err)

	now := time.Now()
	start := now.Add(-24 * time.Hour)
	end := now.Add(24 * time.Hour)
	err = applyDiscUC.Execute(ctx, apply_discount.Request{
		ProductID: productID,
		Percent:   20,
		StartDate: start,
		EndDate:   end,
	})
	require.NoError(t, err)

	product, err := getQuery.Execute(ctx, productID)
	require.NoError(t, err)
	effectiveVal := float64(product.EffectivePriceNumerator) / float64(product.EffectivePriceDenominator)
	assert.Equal(t, 80.0, effectiveVal)
}

func TestActivateDeactivateFlow(t *testing.T) {
	ctx, client, createUC, _, activateUC, deactivateUC, _, _, _, getQuery, listQuery := setup(t)
	productID, err := createUC.Execute(ctx, create_product.Request{
		Name:                 "Active Product",
		Description:          "",
		Category:             "c1",
		BasePriceNumerator:   500,
		BasePriceDenominator: 100,
	})
	require.NoError(t, err)

	product, _ := getQuery.Execute(ctx, productID)
	require.NotNil(t, product)
	assert.Equal(t, "active", product.Status)

	err = deactivateUC.Execute(ctx, deactivate_product.Request{ProductID: productID})
	require.NoError(t, err)
	product, _ = getQuery.Execute(ctx, productID)
	assert.Equal(t, "inactive", product.Status)

	err = activateUC.Execute(ctx, activate_product.Request{ProductID: productID})
	require.NoError(t, err)
	product, _ = getQuery.Execute(ctx, productID)
	assert.Equal(t, "active", product.Status)

	result, err := listQuery.Execute(ctx, contracts.ListFilter{}, contracts.ListPage{PageSize: 10, Token: ""})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(result.Items), 1)
	_ = client
}

func TestBusinessRuleValidation(t *testing.T) {
	ctx, _, createUC, _, _, deactivateUC, _, applyDiscUC, _, _, _ := setup(t)
	productID, err := createUC.Execute(ctx, create_product.Request{
		Name:                 "Inactive",
		Description:          "",
		Category:             "c1",
		BasePriceNumerator:   1000,
		BasePriceDenominator: 100,
	})
	require.NoError(t, err)

	err = deactivateUC.Execute(ctx, deactivate_product.Request{ProductID: productID})
	require.NoError(t, err)

	now := time.Now()
	err = applyDiscUC.Execute(ctx, apply_discount.Request{
		ProductID: productID,
		Percent:   10,
		StartDate: now.Add(-time.Hour),
		EndDate:   now.Add(time.Hour),
	})
	assert.ErrorIs(t, err, domain.ErrProductNotActive)
}

func TestOutboxEventCreation(t *testing.T) {
	ctx, client, createUC, updateUC, activateUC, deactivateUC, _, applyDiscUC, removeDiscUC, _, _ := setup(t)
	productID, err := createUC.Execute(ctx, create_product.Request{
		Name:                 "Outbox Test",
		Description:          "",
		Category:             "c1",
		BasePriceNumerator:   1000,
		BasePriceDenominator: 100,
	})
	require.NoError(t, err)

	events := getOutboxEvents(t, client, productID)
	require.Len(t, events, 1)
	assert.Equal(t, "product.created", events[0].EventType)

	err = updateUC.Execute(ctx, update_product.Request{ProductID: productID, Name: "Updated", Description: "", Category: "c1"})
	require.NoError(t, err)
	events = getOutboxEvents(t, client, productID)
	require.GreaterOrEqual(t, len(events), 2)

	err = deactivateUC.Execute(ctx, deactivate_product.Request{ProductID: productID})
	require.NoError(t, err)
	err = activateUC.Execute(ctx, activate_product.Request{ProductID: productID})
	require.NoError(t, err)

	now := time.Now()
	err = applyDiscUC.Execute(ctx, apply_discount.Request{ProductID: productID, Percent: 5, StartDate: now.Add(-time.Hour), EndDate: now.Add(time.Hour)})
	require.NoError(t, err)
	err = removeDiscUC.Execute(ctx, remove_discount.Request{ProductID: productID})
	require.NoError(t, err)

	events = getOutboxEvents(t, client, productID)
	require.GreaterOrEqual(t, len(events), 5)
}

func TestArchiveFlow(t *testing.T) {
	ctx, client, createUC, _, _, _, archiveUC, _, _, getQuery, listQuery := setup(t)
	productID, err := createUC.Execute(ctx, create_product.Request{
		Name:                 "To Archive",
		Description:          "",
		Category:             "c1",
		BasePriceNumerator:   1000,
		BasePriceDenominator: 100,
	})
	require.NoError(t, err)

	product, err := getQuery.Execute(ctx, productID)
	require.NoError(t, err)
	assert.Equal(t, "active", product.Status)

	err = archiveUC.Execute(ctx, archive_product.Request{ProductID: productID})
	require.NoError(t, err)

	product, err = getQuery.Execute(ctx, productID)
	require.NoError(t, err)
	assert.Equal(t, "archived", product.Status)

	events := getOutboxEvents(t, client, productID)
	require.GreaterOrEqual(t, len(events), 2)
	var foundArchived bool
	for _, e := range events {
		if e.EventType == "product.archived" {
			foundArchived = true
			break
		}
	}
	assert.True(t, foundArchived, "expected product.archived event in outbox")

	// Archived products do not appear in list (active only)
	result, err := listQuery.Execute(ctx, contracts.ListFilter{}, contracts.ListPage{PageSize: 100, Token: ""})
	require.NoError(t, err)
	for _, item := range result.Items {
		assert.NotEqual(t, productID, item.ID, "archived product should not be in list")
	}
}

func TestConcurrentUpdates(t *testing.T) {
	ctx, _, createUC, updateUC, _, _, _, _, _, getQuery, _ := setup(t)
	productID, err := createUC.Execute(ctx, create_product.Request{
		Name:                 "Original",
		Description:          "",
		Category:             "c1",
		BasePriceNumerator:   1000,
		BasePriceDenominator: 100,
	})
	require.NoError(t, err)

	var wg sync.WaitGroup
	var err1, err2 error
	wg.Add(2)
	go func() {
		defer wg.Done()
		err1 = updateUC.Execute(ctx, update_product.Request{
			ProductID: productID, Name: "Update A", Description: "", Category: "c1",
		})
	}()
	go func() {
		defer wg.Done()
		err2 = updateUC.Execute(ctx, update_product.Request{
			ProductID: productID, Name: "Update B", Description: "", Category: "c1",
		})
	}()
	wg.Wait()

	successCount := 0
	if err1 == nil {
		successCount++
	}
	if err2 == nil {
		successCount++
	}
	require.Equal(t, 1, successCount, "exactly one update should succeed")
	if err1 != nil {
		assert.ErrorIs(t, err1, commitplan.ErrConcurrentModification)
	}
	if err2 != nil {
		assert.ErrorIs(t, err2, commitplan.ErrConcurrentModification)
	}

	product, err := getQuery.Execute(ctx, productID)
	require.NoError(t, err)
	require.NotNil(t, product)
	assert.Contains(t, []string{"Update A", "Update B"}, product.Name)
}
