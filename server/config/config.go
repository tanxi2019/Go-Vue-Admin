package config

import "go.uber.org/zap/zapcore"

// 系统配置，对应yml
// viper内置了map structure, yml文件用"-"区分单词, 转为驼峰方便

// Conf 全局配置变量
var Conf = new(config)

type config struct {
	System    *SystemConfig    `mapstructure:"system" json:"system"`
	Logs      *LogsConfig      `mapstructure:"logs" json:"logs"`
	Mysql     *MysqlConfig     `mapstructure:"mysql" json:"mysql"`
	Redis     *RedisConfig     `mapstructure:"redis" json:"redis"`
	Casbin    *CasbinConfig    `mapstructure:"casbin" json:"casbin"`
	Jwt       *JwtConfig       `mapstructure:"jwt" json:"jwt"`
	RateLimit *RateLimitConfig `mapstructure:"rate-limit" json:"rateLimit"`
	Qiniu     *QiniuConfig     `mapstructure:"qiniu" json:"qiniu"`
}

type SystemConfig struct {
	Mode          string `mapstructure:"mode" json:"mode"`
	UrlPathPrefix string `mapstructure:"url-path-prefix" json:"urlPathPrefix"`
	Host          string `mapstructure:"host" json:"host"`
	Port          int    `mapstructure:"port" json:"port"`
	I18nLanguage  string `mapstructure:"i18n-language" json:"i18nLanguage"`
}

type LogsConfig struct {
	Level      zapcore.Level `mapstructure:"level" json:"level"`
	Path       string        `mapstructure:"path" json:"path"`
	MaxSize    int           `mapstructure:"max-size" json:"maxSize"`
	MaxBackups int           `mapstructure:"max-backups" json:"maxBackups"`
	MaxAge     int           `mapstructure:"max-age" json:"maxAge"`
	Compress   bool          `mapstructure:"compress" json:"compress"`
}

type MysqlConfig struct {
	Username  string `mapstructure:"username" json:"username"`
	Password  string `mapstructure:"password" json:"password"`
	Database  string `mapstructure:"database" json:"database"`
	Host      string `mapstructure:"host" json:"host"`
	Port      int    `mapstructure:"port" json:"port"`
	Query     string `mapstructure:"query" json:"query"`
	LogMode   bool   `mapstructure:"log-mode" json:"logMode"`
	Charset   string `mapstructure:"charset" json:"charset"`
	Collation string `mapstructure:"collation" json:"collation"`
}

type RedisConfig struct {
	Password string `mapstructure:"password" json:"password"`
	Database int    `mapstructure:"database" json:"database"`
	Addr     string `mapstructure:"addr" json:"addr"`
}

type CasbinConfig struct {
	ModelPath string `mapstructure:"model-path" json:"modelPath"`
}

type JwtConfig struct {
	Realm      string `mapstructure:"realm" json:"realm"`
	Key        string `mapstructure:"key" json:"key"`
	Timeout    int    `mapstructure:"timeout" json:"timeout"`
	MaxRefresh int    `mapstructure:"max-refresh" json:"maxRefresh"`
}

type RateLimitConfig struct {
	FillInterval int64 `mapstructure:"fill-interval" json:"fillInterval"`
	Capacity     int64 `mapstructure:"capacity" json:"capacity"`
	Quantum      int64 `mapstructure:"quantum" json:"quantum"`
}

type QiniuConfig struct {
	AccessKey string `mapstructure:"accessKey" json:"accessKey"`
	SecretKey string `mapstructure:"secretKey" json:"secretKey"`
	Bucket    string `mapstructure:"bucket" json:"bucket"`
	Origin    string `mapstructure:"origin" json:"origin"`
}
