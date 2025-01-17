package dolt_test

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/dolt"
)

func ExampleRunContainer() {
	// runDoltContainer {
	ctx := context.Background()

	doltContainer, err := dolt.RunContainer(ctx,
		testcontainers.WithImage("dolthub/dolt-sql-server:1.32.4"),
		dolt.WithConfigFile(filepath.Join("testdata", "dolt_1_32_4.cnf")),
		dolt.WithDatabase("foo"),
		dolt.WithUsername("root"),
		dolt.WithPassword("password"),
		dolt.WithScripts(filepath.Join("testdata", "schema.sql")),
	)
	if err != nil {
		panic(err)
	}

	// Clean up the container
	defer func() {
		if err := doltContainer.Terminate(ctx); err != nil {
			panic(err)
		}
	}()
	// }

	state, err := doltContainer.State(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(state.Running)

	// Output:
	// true
}

func ExampleRunContainer_connect() {
	ctx := context.Background()

	doltContainer, err := dolt.RunContainer(ctx,
		testcontainers.WithImage("dolthub/dolt-sql-server:1.32.4"),
		dolt.WithConfigFile(filepath.Join("testdata", "dolt_1_32_4.cnf")),
		dolt.WithDatabase("foo"),
		dolt.WithUsername("bar"),
		dolt.WithPassword("password"),
		dolt.WithScripts(filepath.Join("testdata", "schema.sql")),
	)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := doltContainer.Terminate(ctx); err != nil {
			panic(err)
		}
	}()

	connectionString, _ := doltContainer.ConnectionString(ctx)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err)
	}
	stmt, err := db.Prepare("SELECT dolt_version();")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	row := stmt.QueryRow()
	version := ""
	err = row.Scan(&version)
	if err != nil {
		panic(err)
	}

	fmt.Println(version)

	// Output:
	// 1.32.4
}
