package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/go-playground/validator"
	"github.com/go-yaml/yaml"
	"github.com/hashicorp/vault/api"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	configFlag = flag.String("config", "", "path to config file")
	logger     *zap.Logger
	config     *Config
	vault      *api.Logical
	err        error
	conn       *sql.DB
	initDBDone chan bool
	client     *twitch.Client
)

func main() {
	flag.Parse()
	logger, _ = zap.NewDevelopment()

	config, err = NewYAMLConfig(*configFlag)
	if err != nil {
		logger.Fatal("Can not load yaml config", zap.String("path", *configFlag), zap.String("error", err.Error()))
	}
	logger.Info("Bot config loaded", zap.String("config", config.Dump()))

	var addr = "https://192.168.122.36:8200"
	var token = "s.J3ZN4w0fb5ug3WRR61i652eZ"
	vault, err = NewVault(addr, token)
	if err != nil {
		logger.Fatal("Can not create connection to vault", zap.String("address", addr), zap.String("token", token), zap.String("error", err.Error()))
	}

	initDBDone = make(chan bool)
	go func() {
		conn, err = NewDBConn(vault, config)
		if err != nil {
			logger.Fatal("Can not create database connection", zap.String("config", config.Dump()), zap.String("error", err.Error()))
		}
		initDBDone <- true
		logger.Info("DB connection created")
	}()

	client, err = NewTwitchClient(vault, config)
	if err != nil {
		logger.Fatal("Can not create twitch client", zap.String("config", config.Dump()), zap.String("error", err.Error()))
	}
	logger.Info("Twitch client created")

	var queryStr = "CALL add_data($1, $2, $3, $4, $5, $6, $7, $8)"
	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		var role = msg.Tags["turbo"] + msg.Tags["mod"] + msg.Tags["subscriber"]

		_, err = conn.Exec(queryStr, msg.ID, msg.Time, msg.RoomID, msg.Channel,
			msg.User.ID, msg.User.DisplayName, role, msg.Message)
		if err != nil {
			logger.Warn("Can not add data to database", zap.String("error", err.Error()))
		}
	})

	<-initDBDone
	logger.Info("Start connecting to twitch")
	err = client.Connect()
	if err != nil {
		logger.Warn("Error during connection to twitch", zap.String("error", err.Error()))
	}
}

func NewDBConn(v *api.Logical, c *Config) (*sql.DB, error) {
	secret, err := v.Read("twchd/postgres")
	if err != nil {
		return nil, err
	}
	var user = secret.Data["user"].(string)
	var passwd = secret.Data["password"].(string)
	var connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Address, c.Port, user, passwd, user)
	return sql.Open("postgres", connStr)
}

func NewTwitchClient(v *api.Logical, c *Config) (*twitch.Client, error) {
	secret, err := v.Read("twchd/twitch")
	if err != nil {
		return nil, err
	}
	var user = secret.Data["username"].(string)
	var token = secret.Data["token"].(string)

	var client = twitch.NewClient(user, token)
	client.TLS = false
	for _, channel := range c.AccountsList {
		client.Join(channel)
	}
	return client, nil
}

type Config struct {
	AccountsList []string `yaml:"join_to" validate:"required,unique"`
	Address      string   `yaml:"address" validate:"ipv4"`
	Port         int      `yaml:"port" validate:"min=1025,max=65535"`
}

func NewYAMLConfig(filename string) (*Config, error) {
	var validate = validator.New()
	var err = validate.Var(filename, "file,endswith=.yml|endswith=.yaml")
	if err != nil {
		return nil, err
	}
	rawConfig, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config = new(Config)
	err = yaml.Unmarshal(rawConfig, config)
	if err != nil {
		return nil, err
	}
	err = validate.Struct(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) Dump() string {
	return fmt.Sprintf("%+v\n", c)
}

func NewVault(url, token string) (*api.Logical, error) {
	assetCrt, err := assets.Open("/vault.crt")
	if err != nil {
		return nil, err
	}
	pem, err := ioutil.ReadAll(assetCrt)
	if err != nil {
		return nil, err
	}
	var certPool = x509.NewCertPool()
	var ok = certPool.AppendCertsFromPEM(pem)
	if !ok {
		return nil, errors.New("Can not load pem crt from vfs")
	}
	var config = api.DefaultConfig()
	config.HttpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{RootCAs: certPool}
	vault, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	err = vault.SetAddress(url)
	if err != nil {
		return nil, err
	}
	vault.SetToken(token)
	return vault.Logical(), nil
}
