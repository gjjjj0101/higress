package embedding

import (
	"encoding/json"
	"errors"

	"github.com/alibaba/higress/plugins/wasm-go/pkg/wrapper"
)

const (
	qwenDomain             = "dashscope.aliyuncs.com"
	qwenChatCompletionPath = "/api/v1/services/embeddings/text-embedding/text-embedding"
)

type dashScopeInitializer struct {
}

func (d *dashScopeInitializer) ValidateConfig(config ProviderConfig) error {
	if config.ApiToken == "" {
		return errors.New("dashscope apiToken is required")
	}
	return nil
}

func (d *dashScopeInitializer) CreateProvider(config ProviderConfig, log wrapper.Log) (Provider, error) {
	log.Debugf("dashscopeInitializer.CreateProvider: %v", config)
	if err := d.ValidateConfig(config); err != nil {
		return nil, err
	}
	s := &dashScopeProvider{}
	log.Debugf("create dashscope embedding client: %v", s)
	s.EmbeddingClient = wrapper.NewClusterClient(wrapper.DnsCluster{
		ServiceName: "dashscope",
		Port:        443,
		Domain:      qwenDomain,
	})
	log.Debugf("dashscope embedding client: %v", s.EmbeddingClient)
	if config.Typ == "" {
		config.Typ = "text-embedding-v1"
	}
	s.Config = config
	log.Debugf("dashscopeInitializer.CreateProvider: success")
	return s, nil
}

type dashScopeProvider struct {
	Config          ProviderConfig
	EmbeddingClient wrapper.HttpClient
}

func (d dashScopeProvider) EmbeddingRequest(raw string, callback wrapper.ResponseCallback) error {
	embeddingRequest := Request{
		Model: d.Config.Typ,
		Input: Input{
			Texts: []string{raw},
		},
		Parameter: Parameter{
			TextType: "query",
		},
	}
	headers := [][2]string{{"Content-Type", "application/json"}, {"Authorization", "Bearer " + d.Config.ApiToken}}
	embeddingRequestSerialized, _ := json.Marshal(embeddingRequest)
	err := d.EmbeddingClient.Post(
		qwenChatCompletionPath,
		headers,
		embeddingRequestSerialized,
		callback)
	if err != nil {
		return err
	}
	return nil
}

func (d dashScopeProvider) ParseEmbeddingResponse(responseBody []byte) ([]float32, error) {
	var responseEmbedding Response
	err := json.Unmarshal(responseBody, &responseEmbedding)
	if err != nil {
		return nil, err
	}
	return responseEmbedding.Output.Embeddings[0].Embedding, nil
}

// DashScope embedding service: Request
type Request struct {
	Model     string    `json:"model"`
	Input     Input     `json:"input"`
	Parameter Parameter `json:"parameters"`
}

type Input struct {
	Texts []string `json:"texts"`
}

type Parameter struct {
	TextType string `json:"text_type"`
}

// DashScope embedding service: Response
type Response struct {
	Output    Output `json:"output"`
	Usage     Usage  `json:"usage"`
	RequestID string `json:"request_id"`
}

type Output struct {
	Embeddings []Embedding `json:"embeddings"`
}

type Embedding struct {
	Embedding []float32 `json:"embedding"`
	TextIndex int32     `json:"text_index"`
}

type Usage struct {
	TotalTokens int32 `json:"total_tokens"`
}
