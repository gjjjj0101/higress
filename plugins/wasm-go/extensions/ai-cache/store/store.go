package store

import (
	"errors"
	"github.com/alibaba/higress/plugins/wasm-go/pkg/wrapper"
)

type StoreConfig struct {
	// @Title zh-CN redis 服务名称
	// @Description zh-CN 带服务类型的完整 FQDN 名称，例如 my-redis.dns、redis.my-ns.svc.cluster.local
	ServiceName string `required:"true" yaml:"serviceName" json:"serviceName"`
	// @Title zh-CN redis 服务端口
	// @Description zh-CN 默认值为6379
	ServicePort int `required:"false" yaml:"servicePort" json:"servicePort"`
	// @Title zh-CN 用户名
	// @Description zh-CN 登陆 redis 的用户名，非必填
	Username string `required:"false" yaml:"username" json:"username"`
	// @Title zh-CN 密码
	// @Description zh-CN 登陆 redis 的密码，非必填，可以只填密码
	Password string `required:"false" yaml:"password" json:"password"`
	// @Title zh-CN 请求超时
	// @Description zh-CN 请求 redis 的超时时间，单位为毫秒。默认值是1000，即1秒
	Timeout int    `required:"false" yaml:"timeout" json:"timeout"`
	Type    string `required:"true" yaml:"type" json:"type"`
	Domain  string `required:"true" yaml:"domain" json:"domain"`
}

const (
	storeTypeTair = "Tair"
)

var (
	errUnsupportedApiName = errors.New("unsupported Store name")

	storeInitializers = map[string]storeInitializer{
		storeTypeTair: &tairInitializer{},
	}
)

type storeInitializer interface {
	ValidateConfig(config StoreConfig) error
	CreateStore(config StoreConfig) (Store, error)
}

type Store interface {
	GetIndex(name string) (Index, error)
	CreateIndex(index Index) (bool, error)
	StoreVector(vector Vector, indexName string) (bool, error)
	SearchVector(index string, topK int, vector []float32, callback wrapper.RedisResponseCallback) error
	ParseSearchResponse(response interface{}) (SearchResponse, error)
	GetKey(index string, key string, callback wrapper.RedisResponseCallback) error
	ParseGetKeyResponse(response interface{}) (KNNSearchResponse, error)
}

func CreateStore(config StoreConfig) (Store, error) {
	if initializer, ok := storeInitializers[config.Type]; ok {
		return initializer.CreateStore(config)
	}
	return nil, errUnsupportedApiName
}
