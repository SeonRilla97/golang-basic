package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var cfg *Config

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Pagination PaginationConfig
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

type PaginationConfig struct {
	DefaultSize int `mapstructure:"default_size"`
	MaxSize     int `mapstructure:"max_size"`
}

// DSN 데이터베이스 연결 문자열 생성
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// Load 설정 파일 로드
func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// 환경 변수 연동
	viper.AutomaticEnv()
	//viper.SetEnvPrefix("GOBOARD") // GOBOARD_DATABASE_PASSWORD 형식으로 사용 민감 정보 환경변수화

	// 환경 변수 키 매핑
	//viper.BindEnv("database.password", "GOBOARD_DB_PASSWORD") // 민감 정보 환경변수화
	//viper.BindEnv("database.host", "GOBOARD_DB_HOST") // 민감 정보 환경변수화

	// 설정 파일 읽기
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("설정 파일 읽기 실패: %w", err)
	}

	// 구조체로 언마샬
	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("설정 파싱 실패: %w", err)
	}

	log.Printf("설정 로드 완료: %s", path)
	return cfg, nil
}

// Get 전역 설정 반환
func Get() *Config {
	return cfg
}
