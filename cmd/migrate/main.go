package main

import (
	"context"
	"os"
	"strings"

	dbadmin "cloud.google.com/go/spanner/admin/database/apiv1"
	databasepb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/api/option"

	"github.com/dawitel/product-catalog-service/internal/pkg/config"
	"github.com/dawitel/product-catalog-service/internal/pkg/logger"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadFromEnv()
	dbPath := cfg.DatabasePath()

	ddl, err := os.ReadFile(cfg.MigrationsPath)
	if err != nil {
		logger.Error("read migrations failed", "path", cfg.MigrationsPath, "err", err)
		os.Exit(1)
	}

	statements := parseDDL(string(ddl))
	adminClient, err := dbadmin.NewDatabaseAdminClient(ctx, option.WithoutAuthentication())
	if err != nil {
		logger.Error("admin client failed", "err", err)
		os.Exit(1)
	}
	defer adminClient.Close()

	op, err := adminClient.UpdateDatabaseDdl(ctx, &databasepb.UpdateDatabaseDdlRequest{
		Database:   dbPath,
		Statements: statements,
	})
	if err != nil {
		logger.Error("update DDL failed", "err", err)
		os.Exit(1)
	}

	if err := op.Wait(ctx); err != nil {
		logger.Error("DDL wait failed", "err", err)
		os.Exit(1)
	}

	logger.Info("migrations applied", "database", dbPath)
}

func parseDDL(s string) []string {
	var stmts []string
	for _, part := range strings.Split(s, ";") {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		stmts = append(stmts, trimmed)
	}
	return stmts
}
