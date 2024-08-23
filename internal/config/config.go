package config

import (
	"strings"
	"sync"

	"ch-common-package/cache"
	"ch-common-package/logger"
	"ch-common-package/mongodb"
	"ch-common-package/ssdb"

	"github.com/BurntSushi/toml"
)

func MustLoad(configPath string) {
	once.Do(func() {
		if _, err := toml.DecodeFile(configPath, &Global); err != nil {
			panic(err)
		}

		if Global.Chain.GenesisTimestamp == 0 {
			//testnet genesis time
			Global.Chain.GenesisTimestamp = 1722964875
		}
	})
}

var (
	Global config
	once   sync.Once
)

type config struct {
	// 运行模式(debug:调试,test:测试,release:正式)
	RunMode string
	HTTP    http
	//是否启用swagger
	Swagger         bool
	SwaggerDocPath  string
	MongoDB         mongodb.Option
	Redis           cache.RedisOption
	LeoCompile      leoCompile `toml:"leoCompile"`
	JwtConfig       Jwt
	WalletRedis     cache.RedisOption
	PushSendMessage bool
	Log             logger.Option
	SSDB            ssdb.Option
	ThirdParty      ThirdParty
	Chain           chainOption
	Nsq             nsqOption
}

type leoCompile struct {
	HomePath string `toml:"homePath"`
	BinPath  string `toml:"binPath"`
}

type http struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

type ThirdParty struct {
	GateioPriceUrl   string
	WeChatNotifyUrl  string
	HeightWarningUrl string
}

type chainOption struct {
	GenesisTimestamp int64
	GenesisHeight    int64
	MainNetTimestamp int64
	MainNetHeight    int64
	//current network, 0: testnet, 1: mainnet
	CurrentNetWork int
}

func (c *config) IsDebugMode() bool {
	return strings.EqualFold(c.RunMode, "debug")
}

type nsqOption struct {
	Lookups  []string
	Nsqds    []string
	ClientId string
	// The server-side message timeout for messages delivered to this client
	Timeout int
}

type Email struct {
	MailDriver      string `toml:"mailDriver"`
	MailHost        string `toml:"mailHost"`
	MailPort        int    `toml:"mailPort"`
	MailUserName    string `toml:"mailUserName"`
	MailFromAddress string `toml:"mailFromAddress"`
	MailPassword    string
	MailEncryption  string
	MailFormName    string
}

type Jwt struct {
	Secret   string
	Timeout  int
	OverTime int
}
