package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/pkg/config"
	"github.com/willie-lin/cloud-terminal/pkg/logger"
)

// findDatabasePath 向上递归搜索非空的 SQLite 数据库文件，以支持在不同工作目录下启动服务时均能连上根目录的数据库。
func findDatabasePath(dbPath string) string {
	if filepath.IsAbs(dbPath) {
		return dbPath
	}
	dir, err := os.Getwd()
	if err != nil {
		return dbPath
	}
	for {
		target := filepath.Join(dir, dbPath)
		if info, err := os.Stat(target); err == nil && !info.IsDir() && info.Size() > 0 {
			return target
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	absPath, err := filepath.Abs(dbPath)
	if err == nil {
		return absPath
	}
	return dbPath
}

// NewClient 创建 Entgo 客户端，配置数据库连接池
func NewClient(dbCfg *config.DatabaseConfig) (*ent.Client, error) {
	var dsn string

	switch dbCfg.Driver {
	case "sqlite", "sqlite3":
		dbCfg.Database = findDatabasePath(dbCfg.Database)

		// Resolve absolute path to avoid confusion
		absPath, err := filepath.Abs(dbCfg.Database)
		if err == nil {
			dbCfg.Database = absPath
		}

		// 🔧 SQLite Optimization: Use WAL mode for better concurrency and set memory cache/temp_store
		// WAL (Write-Ahead Logging) allows concurrent readers while a writer is active
		// _cache_size=-262144 sets 256MB cache, _temp_store=2 uses in-memory temporary tables
		dsn = fmt.Sprintf("%s?_busy_timeout=%d&_fk=1&_journal_mode=WAL&_cache_size=-262144", dbCfg.Database, dbCfg.BusyTimeout)

		// Force SQLite connection pool limits to 1 to prevent self-deadlocks and database locking during migrations
		if dbCfg.MaxOpenConns == 0 {
			dbCfg.MaxOpenConns = 10 // 如果 config 里没写，给个默认兜底值
		}
		if dbCfg.MaxIdleConns == 0 {
			dbCfg.MaxIdleConns = 5 // 如果 config 里没写，给个默认兜底值
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

	logger.Info("Connecting to database",
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

	logger.Info("Opening Ent DB", zap.String("dialect", dialectName))

	drv := entsql.OpenDB(dialectName, sqlDB)
	client := ent.NewClient(ent.Driver(drv))

	logger.Info("Successfully connected to database",
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
	if err := client.Schema.Create(ctx); err != nil {
		logger.Error("Failed creating schema resources", zap.Error(err))
		return fmt.Errorf("failed creating schema resources: %w", err)
	}
	logger.Info("Schema migration completed successfully")
	return nil
}

// DebugMode 调试模式迁移
// 注意：migrate 包需要从 entgo.io/ent/dialect/sql/schema 导入，但在 ent 中，通常是 ent/migrate
func DebugMode(ctx context.Context, client *ent.Client) error {
	// migrate options are usually in the generated ent/migrate package, but here we use the client directly
	// For simple specific debug options, we might need to import the generated migrate package
	// But client.Debug().Schema.Create interface is simpler.
	// To use WithDropIndex and WithDropColumn, we need "github.com/willie-lin/AIStockTrader/ent/migrate"

	// Since we haven't imported the specific migrate package yet, let's keep it simple for now
	// or use standard Create options.

	// Assuming standard Create for now to avoid import cycle or complex setup before generated code is fully ready
	if err := client.Debug().Schema.Create(ctx, schema.WithDropIndex(true), schema.WithDropColumn(true)); err != nil {
		logger.Error("Failed creating schema resources in debug mode", zap.Error(err))
		return fmt.Errorf("failed creating schema resources: %w", err)
	}
	logger.Info("Debug mode schema migration completed successfully")
	return nil
}
