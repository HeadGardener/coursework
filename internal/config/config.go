package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig     DBConfig
	ServerConfig ServerConfig
	RedisConfig  RedisConfig
	TokensConfig TokensConfig
}

type DBConfig struct {
	URL string
}

type ServerConfig struct {
	Port string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type TokensConfig struct {
	SecretKey      string
	AccessTokenTTL time.Duration

	InitialLen      int
	RefreshTokenTTL time.Duration
}

func Init(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		log.Println("invalid config path: ", err)
	}

	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		return nil, errors.New("db name is empty")
	}

	srvport := os.Getenv("SERVER_PORT")
	if srvport == "" {
		return nil, errors.New("server port is empty")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		return nil, errors.New("redis addr is empty")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, errors.New("redis db is empty")
	}

	accessTokenSecretKey := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	if accessTokenSecretKey == "" {
		return nil, errors.New("access token secret key is empty")
	}

	accessTokenTTL, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil {
		return nil, fmt.Errorf("invalid access token ttl: %w", err)
	}

	refreshInitialLen, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_INITIAL_LEN"))
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token initial len: %w", err)
	}

	refreshTokenTTL, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TTL"))
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token ttl: %w", err)
	}

	return &Config{
		DBConfig: DBConfig{
			URL: dburl,
		},
		ServerConfig: ServerConfig{
			Port: srvport,
		},
		RedisConfig: RedisConfig{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       redisDB,
		},
		TokensConfig: TokensConfig{
			SecretKey:       accessTokenSecretKey,
			AccessTokenTTL:  time.Duration(accessTokenTTL) * time.Minute,
			InitialLen:      refreshInitialLen,
			RefreshTokenTTL: time.Duration(refreshTokenTTL) * time.Minute,
		},
	}, nil
}
