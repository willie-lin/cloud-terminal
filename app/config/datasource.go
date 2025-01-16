package config

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/migrate"
	"go.uber.org/zap"
	"os"
	"strconv"
)

type DatabaseCfg struct {
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Type     string `json:"type"`
}

func loadConfigFromFile(path string) (*DatabaseCfg, error) {
	// 使用 os.ReadFile 读取配置文件
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg DatabaseCfg
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &cfg, nil
}

func loadConfigFromEnv() *DatabaseCfg {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		port = 3306
	}

	return &DatabaseCfg{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		Type:     os.Getenv("DB_TYPE"),
	}
}

func NewClient() (*ent.Client, error) {
	var cfg *DatabaseCfg
	var err error

	if _, err = os.Stat("app/config/config.json"); err == nil {
		// 配置文件存在，从文件中读取配置
		cfg, err = loadConfigFromFile("app/config/config.json")
		fmt.Println(cfg)
		if err != nil {
			return nil, err
		}
	} else {
		// 配置文件不存在，从环境变量中读取配置
		cfg = loadConfigFromEnv()
	}

	// 确保 cfg 变量被使用
	var client *ent.Client
	var dataSourceName string

	switch cfg.Type {
	case "sqlite3":
		dataSourceName = fmt.Sprintf("file:%s?_busy_timeout=100000&_fk=1", cfg.DbName)
	case "mysql":
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
	case "postgres", "postgresql":
		dataSourceName = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.DbName, cfg.Password)
	default:
		return client, fmt.Errorf("unknown database type")
	}

	client, err = ent.Open(cfg.Type, dataSourceName)
	if err != nil {
		return client, fmt.Errorf("failed opening connection to %s: %v", cfg.Type, err)
	}

	return client, err
}

func AutoMigration(client *ent.Client, ctx context.Context) {
	log, _ := zap.NewDevelopment()

	if err := client.Schema.Create(ctx); err != nil {
		log.Fatal("failed creating schema resources: %v", zap.Error(err))
		//log.Fatalf("failed creating schema resources: %v", err)
	}
}

func DebugMode(err error, client *ent.Client, ctx context.Context) {
	log, _ := zap.NewDevelopment()

	err = client.Debug().Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)
	if err != nil {
		log.Fatal("failed creating schema resources: %v", zap.Error(err))
		//log.Fatalf("failed creating schema resources: %v", err)
	}
}
