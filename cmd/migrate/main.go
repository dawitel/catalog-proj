package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	dbadmin "cloud.google.com/go/spanner/admin/database/apiv1"
	databasepb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	instanceadmin "cloud.google.com/go/spanner/admin/instance/apiv1"
	instancepb "cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dawitel/catalog-proj/internal/pkg/config"
	"github.com/dawitel/catalog-proj/internal/pkg/logger"
)

const emulatorConfigID = "emulator-config"

func main() {
	if err := run(); err != nil {
		logger.Error("migrate failed", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	cfg := config.LoadFromEnv()

	statements, err := loadStatements(cfg.MigrationsDir)
	if err != nil {
		return err
	}

	opts := []option.ClientOption{option.WithoutAuthentication()}
	dbAdminClient, err := dbadmin.NewDatabaseAdminClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("database admin client: %w", err)
	}
	defer dbAdminClient.Close()

	if cfg.SpannerEmulatorHost != "" {
		return ensureEmulatorInstanceAndDatabase(ctx, cfg, statements, opts, dbAdminClient)
	}

	return applyDDL(ctx, dbAdminClient, cfg.DatabasePath(), statements)
}

func loadStatements(migrationsDir string) ([]string, error) {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("read migrations dir %s: %w", migrationsDir, err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	var statements []string
	for _, name := range names {
		path := filepath.Join(migrationsDir, name)
		ddl, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}
		statements = append(statements, parseDDL(string(ddl))...)
	}
	return statements, nil
}

func applyDDL(ctx context.Context, client *dbadmin.DatabaseAdminClient, database string, statements []string) error {
	op, err := client.UpdateDatabaseDdl(ctx, &databasepb.UpdateDatabaseDdlRequest{
		Database:   database,
		Statements: statements,
	})
	if err != nil {
		return fmt.Errorf("update DDL: %w", err)
	}
	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("DDL wait: %w", err)
	}
	logger.Info("migrations applied", "database", database)
	return nil
}

func ensureEmulatorInstanceAndDatabase(ctx context.Context, cfg config.GlobalConfig, statements []string, opts []option.ClientOption, dbAdminClient *dbadmin.DatabaseAdminClient) error {
	instAdminClient, err := instanceadmin.NewInstanceAdminClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("instance admin client: %w", err)
	}
	defer instAdminClient.Close()

	projectPath := "projects/" + cfg.SpannerProject
	instancePath := projectPath + "/instances/" + cfg.SpannerInstance

	if err := ensureEmulatorInstance(ctx, instAdminClient, projectPath, instancePath, cfg.SpannerInstance); err != nil {
		return err
	}

	op, err := dbAdminClient.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          instancePath,
		CreateStatement: "CREATE DATABASE `" + cfg.SpannerDatabase + "`",
		ExtraStatements: statements,
	})
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			logger.Info("database already exists", "database", cfg.DatabasePath())
			return nil
		}
		return fmt.Errorf("create database: %w", err)
	}
	if _, err := op.Wait(ctx); err != nil {
		return fmt.Errorf("create database wait: %w", err)
	}
	logger.Info("created database with schema", "database", cfg.DatabasePath())
	return nil
}

func ensureEmulatorInstance(ctx context.Context, client *instanceadmin.InstanceAdminClient, projectPath, instancePath, instanceID string) error {
	_, err := client.GetInstance(ctx, &instancepb.GetInstanceRequest{Name: instancePath})
	if err == nil {
		return nil
	}
	if status.Code(err) != codes.NotFound {
		return fmt.Errorf("get instance: %w", err)
	}
	configName := projectPath + "/instanceConfigs/" + emulatorConfigID
	op, err := client.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     projectPath,
		InstanceId: instanceID,
		Instance: &instancepb.Instance{
			Config:      configName,
			DisplayName: instanceID,
			NodeCount:   1,
		},
	})
	if err != nil {
		return fmt.Errorf("create instance: %w", err)
	}
	if _, err := op.Wait(ctx); err != nil {
		return fmt.Errorf("create instance wait: %w", err)
	}
	logger.Info("created emulator instance", "instance", instancePath)
	return nil
}

func parseDDL(s string) []string {
	var stmts []string
	for _, part := range strings.Split(s, ";") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			stmts = append(stmts, trimmed)
		}
	}
	return stmts
}
