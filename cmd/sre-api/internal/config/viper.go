package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds all runtime configuration for Sentinel AI.
type Config struct {
	APIPort string
	APIHost string
	APIEnv  string

	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DatabaseURL string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	KubeConfig  string
	KindCluster string

	PrometheusURL string
	LokiURL       string

	GeminiAPIKey string
	GeminiModel  string

	OrchestratorPort string
	HumanApproval    bool
}

// Load reads configuration from environment variables, .env file, and flags.
// Priority (high to low): flags > env vars > .env file > defaults.
func Load() *Config {
	v := viper.New()

	// --- defaults ---
	v.SetDefault("SRE_API_PORT", "8080")
	v.SetDefault("SRE_API_HOST", "0.0.0.0")
	v.SetDefault("SRE_API_ENV", "development")

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_USER", "sentinel")
	v.SetDefault("DB_PASSWORD", "sentinel")
	v.SetDefault("DB_NAME", "sentinel")

	v.SetDefault("REDIS_ADDR", "localhost:6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	v.SetDefault("KUBECONFIG", defaultKubeConfig())
	v.SetDefault("KIND_CLUSTER_NAME", "sentinel")

	v.SetDefault("PROMETHEUS_URL", "http://localhost:9090")
	v.SetDefault("LOKI_URL", "http://localhost:3100")

	v.SetDefault("LLM_PROVIDER", "gemini")
	v.SetDefault("GEMINI_API_KEY", "")
	v.SetDefault("GEMINI_MODEL", "gemini-1.5-flash")

	v.SetDefault("ORCHESTRATOR_PORT", "8000")
	v.SetDefault("ORCHESTRATOR_ENV", "development")
	v.SetDefault("HUMAN_APPROVAL_REQUIRED", true)

	// --- .env file support ---
	v.AddConfigPath(".")
	v.SetConfigName(".env")
	v.SetConfigType("env")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "warning: failed to read .env file: %v\n", err)
		}
	}

	// --- environment variables ---
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// --- command-line flags ---
	fs := pflag.NewFlagSet("SRE-AI", pflag.ContinueOnError)

	fs.String("api-port", "", "SRE API port")
	fs.String("api-host", "", "SRE API host")
	fs.String("api-env", "", "SRE API environment (development|production)")

	fs.String("db-host", "", "Postgres host")
	fs.String("db-port", "", "Postgres port")
	fs.String("db-user", "", "Postgres user")
	fs.String("db-password", "", "Postgres password")
	fs.String("db-name", "", "Postgres database name")

	fs.String("redis-addr", "", "Redis address")
	fs.String("redis-password", "", "Redis password")
	fs.Int("redis-db", 0, "Redis database index")

	fs.String("kubeconfig", "", "Path to kubeconfig file")
	fs.String("kind-cluster", "", "Kind cluster name")

	fs.String("prometheus-url", "", "Prometheus base URL")
	fs.String("loki-url", "", "Loki base URL")

	fs.String("gemini-api-key", "", "Gemini API key")
	fs.String("gemini-model", "", "Gemini model name")

	fs.String("orchestrator-port", "", "Python orchestrator port")
	fs.Bool("human-approval", false, "Require human approval for destructive remediation")

	// Parse flags; on error, print usage and exit. pflag prints usage itself for --help.
	if err := fs.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "flag parse error: %v\n", err)
		fs.PrintDefaults()
		os.Exit(2)
	}

	// Bind flags to Viper so flag values override env/.env defaults.
	flagBinds := []struct {
		vKey string
		flag *pflag.Flag
	}{
		{"SRE_API_PORT", fs.Lookup("api-port")},
		{"SRE_API_HOST", fs.Lookup("api-host")},
		{"SRE_API_ENV", fs.Lookup("api-env")},
		{"DB_HOST", fs.Lookup("db-host")},
		{"DB_PORT", fs.Lookup("db-port")},
		{"DB_USER", fs.Lookup("db-user")},
		{"DB_PASSWORD", fs.Lookup("db-password")},
		{"DB_NAME", fs.Lookup("db-name")},
		{"REDIS_ADDR", fs.Lookup("redis-addr")},
		{"REDIS_PASSWORD", fs.Lookup("redis-password")},
		{"REDIS_DB", fs.Lookup("redis-db")},
		{"KUBECONFIG", fs.Lookup("kubeconfig")},
		{"KIND_CLUSTER_NAME", fs.Lookup("kind-cluster")},
		{"PROMETHEUS_URL", fs.Lookup("prometheus-url")},
		{"LOKI_URL", fs.Lookup("loki-url")},
		{"GEMINI_API_KEY", fs.Lookup("gemini-api-key")},
		{"GEMINI_MODEL", fs.Lookup("gemini-model")},
		{"ORCHESTRATOR_PORT", fs.Lookup("orchestrator-port")},
		{"HUMAN_APPROVAL_REQUIRED", fs.Lookup("human-approval")},
	}

	for _, b := range flagBinds {
		if err := v.BindPFlag(b.vKey, b.flag); err != nil {
			panic(err)
		}
	}

	return &Config{
		APIPort: v.GetString("SRE_API_PORT"),
		APIHost: v.GetString("SRE_API_HOST"),
		APIEnv:  v.GetString("SRE_API_ENV"),

		DBHost:      v.GetString("DB_HOST"),
		DBPort:      v.GetString("DB_PORT"),
		DBUser:      v.GetString("DB_USER"),
		DBPassword:  v.GetString("DB_PASSWORD"),
		DBName:      v.GetString("DB_NAME"),
		DatabaseURL: v.GetString("DATABASE_URL"),

		RedisAddr:     v.GetString("REDIS_ADDR"),
		RedisPassword: v.GetString("REDIS_PASSWORD"),
		RedisDB:       v.GetInt("REDIS_DB"),

		KubeConfig:  v.GetString("KUBECONFIG"),
		KindCluster: v.GetString("KIND_CLUSTER_NAME"),

		PrometheusURL: v.GetString("PROMETHEUS_URL"),
		LokiURL:       v.GetString("LOKI_URL"),

		GeminiAPIKey: v.GetString("GEMINI_API_KEY"),
		GeminiModel:  v.GetString("GEMINI_MODEL"),

		OrchestratorPort: v.GetString("ORCHESTRATOR_PORT"),
		HumanApproval:    v.GetBool("HUMAN_APPROVAL_REQUIRED"),
	}
}

// defaultKubeConfig returns the standard kubeconfig path for the current user.
func defaultKubeConfig() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}

// String returns a safe, redacted summary of the config (secrets hidden).
func (c *Config) String() string {
	return fmt.Sprintf(
		"Config{API:%s:%s env=%s DB=%s Redis=%s KubeConfig=%s Prometheus=%s Loki=%s GeminiModel=%s Orchestrator=%s HumanApproval=%v}",
		c.APIHost, c.APIPort, c.APIEnv,
		c.DBHost+":"+c.DBPort+"/"+c.DBName,
		c.RedisAddr,
		c.KubeConfig,
		c.PrometheusURL, c.LokiURL,
		c.GeminiModel,
		c.OrchestratorPort, c.HumanApproval,
	)
}

// IsDevelopment returns true when APIEnv is "development".
func (c *Config) IsDevelopment() bool {
	return c.APIEnv == "development"
}

// DBConnectURL returns a lib/pq-compatible connection string.
func (c *Config) DBConnectURL() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

// RequiredStartupDuration is a graceful shutdown timeout for the API server.
const RequiredStartupDuration = 10 * time.Second
