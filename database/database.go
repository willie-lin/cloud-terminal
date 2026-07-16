// package database

// import (
// 	"context"
// 	"entgo.io/ent/dialect/sql"
// 	"fmt"
// 	_ "github.com/go-sql-driver/mysql"
// 	_ "github.com/lib/pq"
// 	_ "modernc.org/sqlite"
// 	"github.com/willie-lin/cloud-terminal/ent"
// 	"github.com/willie-lin/cloud-terminal/ent/migrate"
// 	"go.uber.org/zap"
// 	"time"
// )

// //type Client struct {
// //	Ent *ent.Client
// //	CTX context.Context
// //	Log *log.Logger
// //}

// //type Option

// //func NewClient(opts ...Option) *Client {
// //	cfg := config{log: log.Println, hooks: &hooks{}}
// //	cfg.options(opts...)
// //	client := &Client{config: cfg}
// //	client.init()
// //	return client
// //}

// var (
// 	driver = "mysql"
// 	//baseUrl = "root:root1234@tcp(127.0.0.1:3306)/terminal?charset=utf8&parseTime=true"
// 	baseUrl = "root:root1234@tcp(0.0.0.0:3306)/terminal?charset=utf8&parseTime=true"
// )

// func Client() (*ent.Client, error) {
// 	//drv, err := sql.Open("mysql", "root:root1234@tcp(127.0.0.1:3306)/terminal?charset=utf8&parseTime=true")
// 	drv, err := sql.Open(driver, baseUrl)
// 	fmt.Println(drv)
// 	//drv, err := sql.Open("mysql", "root:root1234@tcp(127.0.0.1:3306)/ent")
// 	if err != nil {
// 		return nil, err
// 	}
// 	db := drv.DB()
// 	//db.SetMaxIdleConns(maxIdleConns)
// 	//db.SetConnMaxLifetime(connMaxLifetime)
// 	//db.SetConnMaxLifetime(connMaxLifetime)
// 	//db.SetMaxOpenConns(maxOpenConns)

// 	db.SetMaxIdleConns(10)
// 	db.SetMaxOpenConns(100)
// 	db.SetConnMaxLifetime(time.Hour)
// 	fmt.Println("ccc")
// 	return ent.NewClient(ent.Driver(drv)), nil

// }

// //func Clients() (*ent.Client, error) {
// //	db, err := sql.Open(driverName, dataSourceName)
// //	if err != nil {
// //		return nil, err
// //	}
// //	//db := drv.DB()
// //	db.SetMaxIdleConns(maxIdleConns)
// //	db.SetConnMaxLifetime(connMaxLifetime)
// //	db.SetMaxOpenConns(maxOpenConns)
// //	drv := entsql.OpenDB(driverName, db)
// //	return ent.NewClient(ent.Driver(drv)), nil
// //}

// func AutoMigration(client *ent.Client, ctx context.Context) {
// 	log, _ := zap.NewDevelopment()
// 	if err := client.Schema.Create(ctx); err != nil {
// 		log.Fatal("failed creating schema resources: %v", zap.Error(err))
// 	}
// }

// func DebugMode(err error, client *ent.Client, ctx context.Context) {
// 	log, _ := zap.NewDevelopment()

// 	err = client.Debug().Schema.Create(
// 		ctx,
// 		migrate.WithDropIndex(true),
// 		migrate.WithDropColumn(true),
// 	)
// 	if err != nil {
// 		log.Fatal("failed creating schema resources: %v", zap.Error(err))
// 	}
// }

package database

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"time"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/willie-lin/ThreatPromptForge/pkg/config"
	"github.com/willie-lin/ThreatPromptForge/pkg/logger"
	sqlite "modernc.org/sqlite"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"go.uber.org/zap"

	"github.com/willie-lin/ThreatPromptForge/ent"
	"github.com/willie-lin/ThreatPromptForge/ent/migrate"
	_ "github.com/willie-lin/ThreatPromptForge/ent/runtime"
)

// func init() {
// 	sql.Register("sqlite3", &sqlite.Driver{})
// }
func init() {
    hook := func(conn sqlite.ExecQuerierContext, dsn string) error {
        if strings.Contains(dsn, "_fk=1") || strings.Contains(dsn, "_fk=true") {
            _, err := conn.ExecContext(context.Background(), "PRAGMA foreign_keys = ON;", nil)
            return err
        }
        return nil
    }

    // Register on the default singleton "sqlite" driver
    sqlite.RegisterConnectionHook(hook)

    // Create and register driver for "sqlite3" alias with the same hook for compatibility with Ent ORM
    drv := &sqlite.Driver{}
    drv.RegisterConnectionHook(hook)
    sql.Register("sqlite3", drv)
}

// NewClient 创建 Entgo 客户端，配置数据库连接池
func NewClient(dbCfg *config.DatabaseConfig) (*ent.Client, error) {
	var dsn string

	switch dbCfg.Driver {
	case "sqlite", "sqlite3":
		// Resolve absolute path to avoid confusion
		absPath, err := filepath.Abs(dbCfg.Database)
		if err == nil {
			dbCfg.Database = absPath
		}

		// 🔧 SQLite Optimization: Use WAL mode for better concurrency
		// WAL (Write-Ahead Logging) allows concurrent readers while a writer is active
		dsn = fmt.Sprintf("file://%s?mode=rwc&_busy_timeout=%d&_fk=1&_journal_mode=WAL&cache=shared", dbCfg.Database, dbCfg.BusyTimeout)

		// Allow multiple connections for transaction support
		// With WAL mode, SQLite can handle concurrent reads and a single writer
		if dbCfg.MaxOpenConns == 0 || dbCfg.MaxOpenConns == 1 {
			dbCfg.MaxOpenConns = 5 // Enough for transaction + concurrent reads
		}
		if dbCfg.MaxIdleConns == 0 || dbCfg.MaxIdleConns == 1 {
			dbCfg.MaxIdleConns = 2
		}
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
			dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Database)
		if dbCfg.MaxOpenConns == 0 {
			dbCfg.MaxOpenConns = 50
		}
		if dbCfg.MaxIdleConns == 0 {
			dbCfg.MaxIdleConns = 10
		}
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			dbCfg.Host, dbCfg.Port, dbCfg.Username, dbCfg.Password, dbCfg.Database, dbCfg.SSLMode)
		if dbCfg.ConnectTimeout > 0 {
			dsn += fmt.Sprintf(" connect_timeout=%d", dbCfg.ConnectTimeout)
		}
		if dbCfg.MaxOpenConns == 0 {
			dbCfg.MaxOpenConns = 50
		}
		if dbCfg.MaxIdleConns == 0 {
			dbCfg.MaxIdleConns = 10
		}
	default:
		logger.Error("Unsupported database driver", zap.String("driver", dbCfg.Driver))
		return nil, fmt.Errorf("unsupported database driver: %s", dbCfg.Driver)
	}

	logger.Info("🔌 Connecting to database",
		zap.String("driver", dbCfg.Driver),
		zap.String("dsn", maskSensitiveDSN(dsn)),
	)

	// 1. 打开数据库连接
	driverName := dbCfg.Driver
	if dbCfg.Driver == "sqlite" {
		driverName = "sqlite3"
	}
	sqlDB, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to %s: %w", driverName, err)
	}

	// 2. 配置连接池
	sqlDB.SetMaxOpenConns(dbCfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbCfg.MaxIdleConns)
	if dbCfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(dbCfg.ConnMaxLifetime)
	} else {
		sqlDB.SetConnMaxLifetime(time.Hour)
	}
	if dbCfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(dbCfg.ConnMaxIdleTime)
	} else {
		sqlDB.SetConnMaxIdleTime(30 * time.Minute)
	}

	// 3. 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed pinging database: %w", err)
	}

	// Fix for SQLite: ensure sqlite_sequence table exists and enable WAL mode explicitly
	if dbCfg.Driver == "sqlite3" || dbCfg.Driver == "sqlite" {
		// Explicitly set WAL mode to ensure concurrency
		if _, err := sqlDB.Exec("PRAGMA journal_mode=WAL;"); err != nil {
			logger.Warn("Failed to set WAL mode", zap.Error(err))
		}
		// Synchronous=NORMAL is safe with WAL and much faster
		if _, err := sqlDB.Exec("PRAGMA synchronous=NORMAL;"); err != nil {
			logger.Warn("Failed to set synchronous mode", zap.Error(err))
		}
		// Increase busy timeout at connection level just in case
		if _, err := sqlDB.Exec(fmt.Sprintf("PRAGMA busy_timeout=%d;", dbCfg.BusyTimeout)); err != nil {
			logger.Warn("Failed to set busy timeout", zap.Error(err))
		}
		// 🚀 Storage & Throughput Optimizations: keep temporary tables/indexes in memory & expand cache
		if _, err := sqlDB.Exec("PRAGMA temp_store=MEMORY;"); err != nil {
			logger.Warn("Failed to set temp_store", zap.Error(err))
		}
		if _, err := sqlDB.Exec("PRAGMA cache_size=-64000;"); err != nil { // 64MB page cache
			logger.Warn("Failed to set cache_size", zap.Error(err))
		}
		if _, err := sqlDB.Exec("PRAGMA mmap_size=268435456;"); err != nil { // 256MB memory mapping
			logger.Warn("Failed to set mmap_size", zap.Error(err))
		}

		_, err := sqlDB.Exec("CREATE TABLE IF NOT EXISTS _dummy_init (id INTEGER PRIMARY KEY AUTOINCREMENT);")
		if err != nil {
			logger.Error("Failed to create dummy table", zap.Error(err))
		} else {
			// Insert a row to trigger sqlite_sequence creation
			sqlDB.Exec("INSERT INTO _dummy_init DEFAULT VALUES;")
			// Clean up
			sqlDB.Exec("DROP TABLE _dummy_init;")
		}
	}

	// 4. 用已打开的连接创建 Ent 客户端（避免重复连接池）
	// Use Ent's dialect constants from entgo.io/ent/dialect package
	dialectName := dbCfg.Driver
	if dbCfg.Driver == "sqlite3" || dbCfg.Driver == "sqlite" {
		dialectName = "sqlite3" // Force "sqlite3"
	} else if dbCfg.Driver == "postgresql" {
		dialectName = dialect.Postgres
	} else if dbCfg.Driver == "mysql" {
		dialectName = dialect.MySQL
	}

	logger.Info("📦 Opening Ent DB", zap.String("dialect", dialectName))

	drv := entsql.OpenDB(dialectName, sqlDB)
	client := ent.NewClient(ent.Driver(drv))

	logger.Info("✅ Successfully connected to database",
		zap.String("driver", dbCfg.Driver),
		zap.String("host", dbCfg.Host),
		zap.String("database", dbCfg.Database),
		zap.Int("max_open_conns", dbCfg.MaxOpenConns),
		zap.Int("max_idle_conns", dbCfg.MaxIdleConns),
		zap.Duration("conn_max_lifetime", dbCfg.ConnMaxLifetime),
		zap.Duration("conn_max_idle_time", dbCfg.ConnMaxIdleTime),
	)

	return client, nil
}

// maskSensitiveDSN 脱敏 DSN
func maskSensitiveDSN(dsn string) string {
	if len(dsn) > 20 {
		return dsn[:20] + "****" + dsn[len(dsn)-20:]
	}
	return "****"
}

// AutoMigration 自动迁移
func AutoMigration(ctx context.Context, client *ent.Client) error {
	// if err := client.Schema.Create(ctx); err != nil {
	// 	logger.Error("Failed creating schema resources", zap.Error(err))
	// 	return fmt.Errorf("failed creating schema resources: %w", err)
	// }
	// logger.Info("Schema migration completed successfully")
	// return nil

	// 添加 migrate.WithForeignKeys(false) 关闭物理外键约束
	if err := client.Schema.Create(ctx, migrate.WithForeignKeys(false)); err != nil {
		logger.Error("Failed creating schema resources", zap.Error(err))
		return fmt.Errorf("failed creating schema resources: %w", err)
	}
	logger.Info("Schema migration completed successfully")
	return nil
}

// AtlasMigrate 使用 Atlas 迁移
func AtlasMigrate(ctx context.Context, client *ent.Client) error {
	if err := client.Schema.Create(ctx, migrate.WithGlobalUniqueID(true)); err != nil {
		logger.Error("Failed creating schema resources with Atlas", zap.Error(err))
		return fmt.Errorf("failed creating schema resources: %w", err)
	}
	logger.Info("Atlas schema migration completed successfully")
	return nil
}

// DebugMode 调试模式迁移
func DebugMode(ctx context.Context, client *ent.Client) error {
	if err := client.Debug().Schema.Create(ctx, migrate.WithDropIndex(true), migrate.WithDropColumn(true)); err != nil {
		logger.Error("Failed creating schema resources in debug mode", zap.Error(err))
		return fmt.Errorf("failed creating schema resources: %w", err)
	}
	logger.Info("Debug mode schema migration completed successfully")
	return nil
}

// WithTx helper function to wrap logic in a transaction
func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

