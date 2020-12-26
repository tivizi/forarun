package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var configContext Context

func init() {
	configBytes, err := ioutil.ReadFile("config/app.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configBytes, &configContext)
	if err != nil {
		panic(err)
	}
	validaion()
}

// Context 配置上下文
type Context struct {
	MongoConfig  MongoConfig  `yaml:"mongodb"`
	MinioConfig  MinioConfig  `yaml:"minio"`
	QcloudConfig QcloudConfig `yaml:"qcloud"`
	SiteConfig   SiteConfig   `yaml:"site"`
	SMTPConfig   SMTPConfig   `yaml:"smtp"`
}

// MongoConfig MongoDB配置文件
type MongoConfig struct {
	URI      string
	Username string
	Password string
}

// MinioConfig 对象存储配置
type MinioConfig struct {
	Enabled   bool
	Endpoint  string
	HTTPS     bool
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
}

// QcloudConfig Qcloud配置文件
type QcloudConfig struct {
	SecretID  string `yaml:"secretID"`
	SecretKey string `yaml:"secretKey"`
	EnableCDN bool   `yaml:"enableCDN"`
}

// SiteConfig 站点配置
type SiteConfig struct {
	Domain string
}

// SMTPConfig 邮箱配置
type SMTPConfig struct {
	Host     string
	Port     int
	SSL      bool `yaml:"ssl"`
	Account  string
	Password string
	Enabled  bool
}

// GetContext 获取配置上下文
func GetContext() Context {
	return configContext
}

func validaion() {
	if len(configContext.SiteConfig.Domain) == 0 {
		panic("site.domain: 站点名是必须")
	}
	if configContext.MinioConfig.Enabled {
		if len(configContext.MinioConfig.AccessKey) == 0 {
			panic("minio.accessKey: 启用对象存储后，此项是必须的")
		}
		if len(configContext.MinioConfig.SecretKey) == 0 {
			panic("minio.secretKey: 启用对象存储后，此项是必须的")
		}
	}
}
