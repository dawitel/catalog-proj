package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	dbadmin "cloud.google.com/go/spanner/admin/database/apiv1"
	databasepb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	project := os.Getenv("SPANNER_PROJECT")
	instance := os.Getenv("SPANNER_INSTANCE")
	database := os.Getenv("SPANNER_DATABASE")

	if project == "" {
		project = "test-project"
	}

	if instance == "" {
		instance = "test-instance"
	}

	if database == "" {
		database = "product-catalog"
	}

	dbPath := fmt.Sprintf("projects/%s/instances/%s/databases/%s", project, instance, database)
	ddl, err := os.ReadFile("migrations/001_initial_schema.sql")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read migrations: %v\n", err)
		os.Exit(1)
	}

	statements := parseDDL(string(ddl))
	adminClient, err := dbadmin.NewDatabaseAdminClient(ctx, option.WithoutAuthentication())
	if err != nil {
		fmt.Fprintf(os.Stderr, "admin client: %v\n", err)
		os.Exit(1)
	}
	defer adminClient.Close()

	op, err := adminClient.UpdateDatabaseDdl(ctx, &databasepb.UpdateDatabaseDdlRequest{
		Database:   dbPath,
		Statements: statements,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "update DDL: %v\n", err)
		os.Exit(1)
	}

	if err := op.Wait(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "DDL wait: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("migrations applied")
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
