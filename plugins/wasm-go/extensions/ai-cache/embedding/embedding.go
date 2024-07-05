package embedding

import (
	"errors"

	"github.com/alibaba/higress/plugins/wasm-go/pkg/wrapper"
)

const (
	providerTypeDashScope = "dashscope"
)

var (
	errUnsupportedProviderType = errors.New("unsupported provider type")

	providerInitializers = map[string]providerInitializer{
		providerTypeDashScope: &dashScopeInitializer{},
	}
)

type ProviderConfig struct {
	Typ      string `json:"type"`
	ApiToken string `json:"apiToken"`
	Model    string `json:"model"`
}

type providerInitializer interface {
	ValidateConfig(config ProviderConfig) error
	CreateProvider(config ProviderConfig, log wrapper.Log) (Provider, error)
}

type Provider interface {
	EmbeddingRequest(raw string, callback wrapper.ResponseCallback) error
	ParseEmbeddingResponse(responseBody []byte) ([]float32, error)
}

func CreateProvider(config ProviderConfig, log wrapper.Log) (Provider, error) {
	initializer, has := providerInitializers[config.Model]
	if !has {
		return nil, errUnsupportedProviderType
	}
	return initializer.CreateProvider(config, log)
}
