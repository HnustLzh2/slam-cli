package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const errNotLoggedIn = "用户未登录，请先执行 'slam-cli auth login'"

type Config struct {
	Profile string     `json:"profile" yaml:"profile"`
	Output  string     `json:"output" yaml:"output"`
	Auth    AuthConfig `json:"auth" yaml:"auth"`
}

type AuthConfig struct {
	LoggedIn     bool      `json:"logged_in" yaml:"logged_in"`
	Region       string    `json:"region" yaml:"region"`
	TokenFile    string    `json:"token_file" yaml:"token_file"`
	CookieJar    string    `json:"cookie_jar" yaml:"cookie_jar"`
	TokenPreview string    `json:"token_preview" yaml:"token_preview"`
	User         UserInfo  `json:"user" yaml:"user"`
	UpdatedAt    time.Time `json:"updated_at" yaml:"updated_at"`
}

type UserInfo struct {
	Username    string `json:"username" yaml:"username"`
	Email       string `json:"email" yaml:"email"`
	DisplayName string `json:"display_name" yaml:"display_name"`
	AvatarURL   string `json:"avatar_url" yaml:"avatar_url"`
	Department  string `json:"department" yaml:"department"`
	OpenID      string `json:"open_id" yaml:"open_id"`
}

func (a AuthConfig) IsLoggedIn() bool {
	if !a.LoggedIn {
		return false
	}

	if a.TokenFile == "" {
		return false
	}

	if _, err := os.Stat(a.TokenFile); err != nil {
		return false
	}

	return true
}

func Default() Config {
	return Config{
		Profile: "default",
		Output:  "plain",
		Auth: AuthConfig{
			LoggedIn: false,
			Region:   "cn",
		},
	}
}

func BaseDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve user home directory: %w", err)
	}

	baseDir := filepath.Join(homeDir, ".slam-cli")
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create base directory: %w", err)
	}

	return baseDir, nil
}

func FilePath() (string, error) {
	baseDir, err := BaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, "config.json"), nil
}

func Load() (Config, error) {
	path, err := FilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := Default()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

func Save(cfg Config) error {
	path, err := FilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func SaveAuth(authCfg AuthConfig) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	cfg.Auth = authCfg
	return Save(cfg)
}

func CurrentAuth() (AuthConfig, error) {
	cfg, err := Load()
	if err != nil {
		return AuthConfig{}, err
	}

	if !cfg.Auth.IsLoggedIn() {
		cfg.Auth.LoggedIn = false
	}

	return cfg.Auth, nil
}

func IsLoggedIn() (bool, error) {
	authCfg, err := CurrentAuth()
	if err != nil {
		return false, err
	}
	return authCfg.IsLoggedIn(), nil
}

func RequireLogin() (AuthConfig, error) {
	authCfg, err := CurrentAuth()
	if err != nil {
		return AuthConfig{}, err
	}

	if !authCfg.IsLoggedIn() {
		return AuthConfig{}, fmt.Errorf(errNotLoggedIn)
	}

	return authCfg, nil
}

func MarkLoggedOut(region string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	if region != "" {
		cfg.Auth.Region = region
	}
	cfg.Auth.LoggedIn = false
	cfg.Auth.TokenFile = ""
	cfg.Auth.CookieJar = ""
	cfg.Auth.TokenPreview = ""
	cfg.Auth.User = UserInfo{}
	cfg.Auth.UpdatedAt = time.Now()

	return Save(cfg)
}

func PreviewToken(token string) string {
	if len(token) <= 16 {
		return token
	}
	return token[:8] + "..." + token[len(token)-8:]
}
