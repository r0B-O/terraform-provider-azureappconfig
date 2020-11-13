package pg

// Create the test database: createdb terraform_backend_pg_test
// TF_ACC=1 GO111MODULE=on go test -v -mod=vendor -timeout=2m -parallel=4 github.com/hashicorp/terraform/backend/remote-state/pg

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/backend"
	"github.com/hashicorp/terraform/states/remote"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// Function to skip a test unless in ACCeptance test mode.
//
// A running Postgres server identified by env variable
// DATABASE_URL is required for acceptance tests.
func testACC(t *testing.T) {
	skip := os.Getenv("TF_ACC") == ""
	if skip {
		t.Log("pg backend tests require setting TF_ACC")
		t.Skip()
	}
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", "postgres://localhost/terraform_backend_pg_test?sslmode=disable")
	}
}

func TestBackend_impl(t *testing.T) {
	var _ backend.Backend = new(Backend)
}

func TestBackendConfig(t *testing.T) {
	testACC(t)
	connStr := getDatabaseUrl()
	schemaName := pq.QuoteIdentifier(fmt.Sprintf("terraform_%s", t.Name()))

	config := backend.TestWrapConfig(map[string]interface{}{
		"conn_str":    connStr,
		"schema_name": schemaName,
	})
	schemaName = pq.QuoteIdentifier(schemaName)

	dbCleaner, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer dbCleaner.Query(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schemaName))

	b := backend.TestBackendConfig(t, New(), config).(*Backend)

	if b == nil {
		t.Fatal("Backend could not be configured")
	}

	_, err = b.db.Query(fmt.Sprintf("SELECT name, data FROM %s.%s LIMIT 1", schemaName, statesTableName))
	if err != nil {
		t.Fatal(err)
	}

	_, err = b.StateMgr(backend.DefaultStateName)
	if err != nil {
		t.Fatal(err)
	}

	s, err := b.StateMgr(backend.DefaultStateName)
	if err != nil {
		t.Fatal(err)
	}
	c := s.(*remote.State).Client.(*RemoteClient)
	if c.Name != backend.DefaultStateName {
		t.Fatal("RemoteClient name is not configured")
	}

	backend.TestBackendStates(t, b)
}

func TestBackendConfigSkipOptions(t *testing.T) {
	testACC(t)
	connStr := getDatabaseUrl()

	testCases := []struct {
		Name               string
		SkipSchemaCreation bool
		SkipTableCreation  bool
		SkipIndexCreation  bool
		Setup              func(t *testing.T, db *sql.DB, schemaName string)
	}{
		{
			Name:               "skip_schema_creation",
			SkipSchemaCreation: true,
			Setup: func(t *testing.T, db *sql.DB, schemaName string) {
				// create the schema as a prerequisites
				_, err := db.Query(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, schemaName))
				if err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Name:              "skip_table_creation",
			SkipTableCreation: true,
			Setup: func(t *testing.T, db *sql.DB, schemaName string) {
				// since the table needs to be already created the schema must be too
				_, err := db.Query(fmt.Sprintf(`CREATE SCHEMA %s`, schemaName))
				if err != nil {
					t.Fatal(err)
				}
				_, err = db.Query(fmt.Sprintf(`CREATE TABLE %s.%s (
					id SERIAL PRIMARY KEY,
					name TEXT,
					data TEXT
					)`, schemaName, statesTableName))
				if err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Name:              "skip_index_creation",
			SkipIndexCreation: true,
			Setup: func(t *testing.T, db *sql.DB, schemaName string) {
				// Everything need to exists for the index to be created
				_, err := db.Query(fmt.Sprintf(`CREATE SCHEMA %s`, schemaName))
				if err != nil {
					t.Fatal(err)
				}
				_, err = db.Query(fmt.Sprintf(`CREATE TABLE %s.%s (
					id SERIAL PRIMARY KEY,
					name TEXT,
					data TEXT
					)`, schemaName, statesTableName))
				if err != nil {
					t.Fatal(err)
				}
				_, err = db.Exec(fmt.Sprintf(`CREATE UNIQUE INDEX IF NOT EXISTS %s ON %s.%s (name)`, statesIndexName, schemaName, statesTableName))
				if err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			schemaName := tc.Name

			config := backend.TestWrapConfig(map[string]interface{}{
				"conn_str":             connStr,
				"schema_name":          schemaName,
				"skip_schema_creation": tc.SkipSchemaCreation,
				"skip_table_creation":  tc.SkipTableCreation,
				"skip_index_creation":  tc.SkipIndexCreation,
			})
			schemaName = pq.QuoteIdentifier(schemaName)
			db, err := sql.Open("postgres", connStr)
			if err != nil {
				t.Fatal(err)
			}

			tc.Setup(t, db, schemaName)
			defer db.Query(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schemaName))

			b := backend.TestBackendConfig(t, New(), config).(*Backend)

			if b == nil {
				t.Fatal("Backend could not be configured")
			}

			// Make sure everything has been created

			// This tests that both the schema and the table have been created
			_, err = b.db.Query(fmt.Sprintf("SELECT name, data FROM %s.%s LIMIT 1", schemaName, statesTableName))
			if err != nil {
				t.Fatal(err)
			}
			// Make sure that the index exists
			query := `select count(*) from pg_indexes where schemaname=$1 and tablename=$2 and indexname=$3;`
			var count int
			if err := b.db.QueryRow(query, tc.Name, statesTableName, statesIndexName).Scan(&count); err != nil {
				t.Fatal(err)
			}
			if count != 1 {
				t.Fatalf("The index has not been created (%d)", count)
			}

			_, err = b.StateMgr(backend.DefaultStateName)
			if err != nil {
				t.Fatal(err)
			}

			s, err := b.StateMgr(backend.DefaultStateName)
			if err != nil {
				t.Fatal(err)
			}
			c := s.(*remote.State).Client.(*RemoteClient)
			if c.Name != backend.DefaultStateName {
				t.Fatal("RemoteClient name is not configured")
			}
		})
	}

}

func TestBackendStates(t *testing.T) {
	testACC(t)
	connStr := getDatabaseUrl()

	testCases := []string{
		fmt.Sprintf("terraform_%s", t.Name()),
		fmt.Sprintf("test with spaces: %s", t.Name()),
	}
	for _, schemaName := range testCases {
		t.Run(schemaName, func(t *testing.T) {
			dbCleaner, err := sql.Open("postgres", connStr)
			if err != nil {
				t.Fatal(err)
			}
			defer dbCleaner.Query("DROP SCHEMA IF EXISTS %s CASCADE", pq.QuoteIdentifier(schemaName))

			config := backend.TestWrapConfig(map[string]interface{}{
				"conn_str":    connStr,
				"schema_name": schemaName,
			})
			b := backend.TestBackendConfig(t, New(), config).(*Backend)

			if b == nil {
				t.Fatal("Backend could not be configured")
			}

			backend.TestBackendStates(t, b)
		})
	}
}

func TestBackendStateLocks(t *testing.T) {
	testACC(t)
	connStr := getDatabaseUrl()
	schemaName := fmt.Sprintf("terraform_%s", t.Name())
	dbCleaner, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer dbCleaner.Query(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schemaName))

	config := backend.TestWrapConfig(map[string]interface{}{
		"conn_str":    connStr,
		"schema_name": schemaName,
	})
	b := backend.TestBackendConfig(t, New(), config).(*Backend)

	if b == nil {
		t.Fatal("Backend could not be configured")
	}

	bb := backend.TestBackendConfig(t, New(), config).(*Backend)

	if bb == nil {
		t.Fatal("Backend could not be configured")
	}

	backend.TestBackendStateLocks(t, b, bb)
}

func getDatabaseUrl() string {
	return os.Getenv("DATABASE_URL")
}
