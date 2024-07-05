package store

import (
	"errors"
	"github.com/alibaba/higress/plugins/wasm-go/pkg/wrapper"
	"github.com/google/uuid"
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/tidwall/resp"
)

type tairInitializer struct {
}

func (t *tairInitializer) ValidateConfig(config StoreConfig) error {
	if config.ServiceName == "" {
		return errors.New("serviceName is required")
	}
	return nil
}

func (t *tairInitializer) CreateStore(config StoreConfig) (Store, error) {
	if err := t.ValidateConfig(config); err != nil {
		return nil, err
	}
	s := &TairStore{}
	s.Client = wrapper.NewRedisClusterClient(wrapper.FQDNCluster{
		FQDN: config.ServiceName,
		Port: int64(config.ServicePort),
	})
	if err := s.Client.Init(config.Username, config.Password, int64(config.Timeout)); err != nil {
		return nil, err
	}
	s.Config = config
	return s, nil
}

type TairStore struct {
	Config StoreConfig

	Client wrapper.RedisClient
}

func (t *TairStore) GetIndex(name string) (Index, error) {
	return Index{}, nil
}

func (t *TairStore) CreateIndex(index Index) (bool, error) {
	if err := ValidateIndexConfig(index); err != nil {
		return false, err
	}
	args := make([]interface{}, 0)
	args = append(args, "TVS.CREATEINDEX")
	args = append(args, index.Name)
	args = append(args, index.Dim)
	args = append(args, index.Typ)
	args = append(args, index.DistanceMethod)
	args = append(args, "data_type")
	args = append(args, index.DataType)
	err := t.Client.Command(args, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (t *TairStore) StoreVector(vector Vector, indexName string) (bool, error) {
	key := uuid.New().String()
	args := make([]interface{}, 0)
	args = append(args, "TVS.HSET")
	args = append(args, indexName)
	args = append(args, key)
	args = append(args, "VECTOR", VectorToString(vector))
	args = append(args, "raw", vector.Raw)
	args = append(args, "answer", vector.Answer)
	if err := t.Client.Command(args, nil); err != nil {
		return false, err
	}
	return true, nil
}

func (t *TairStore) SearchVector(index string, topK int, vector []float32, callback wrapper.RedisResponseCallback) error {
	args := make([]interface{}, 0)
	args = append(args, "TVS.KNNSEARCH")
	args = append(args, index)
	args = append(args, topK)
	args = append(args, VectorToString(Vector{Content: vector}))
	proxywasm.LogDebugf("SearchVector: %v", args)
	err := t.Client.Command(args, callback)
	if err != nil {
		return err
	}
	return nil
}

func (t *TairStore) GetKey(index string, key string, callback wrapper.RedisResponseCallback) error {
	args := make([]interface{}, 0)
	args = append(args, "TVS.HMGET")
	args = append(args, index)
	args = append(args, key)
	args = append(args, "raw")
	args = append(args, "answer")
	err := t.Client.Command(args, callback)
	if err != nil {
		return err
	}
	return nil
}

func (t *TairStore) ParseGetKeyResponse(response interface{}) (KNNSearchResponse, error) {
	value, ok := response.(resp.Value)
	if !ok {
		return KNNSearchResponse{}, errors.New("invalid response")
	}
	if err := value.Error(); err != nil {
		return KNNSearchResponse{}, err
	}
	if value.IsNull() {
		return KNNSearchResponse{}, errors.New("cache miss")
	}
	knnSearchResponse := KNNSearchResponse{}
	knnSearchResponse.Key = value.Array()[0].String()
	knnSearchResponse.Distance = float32(value.Array()[1].Float())
	return knnSearchResponse, nil
}

func (t *TairStore) ParseSearchResponse(response interface{}) (SearchResponse, error) {

	value, ok := response.(resp.Value)
	if !ok {
		return SearchResponse{}, errors.New("invalid response")
	}
	if err := value.Error(); err != nil {
		return SearchResponse{}, err
	}
	if value.IsNull() {
		return SearchResponse{}, errors.New("no such field")
	}
	res := SearchResponse{}
	res.Raw = value.Array()[0].String()
	res.Answer = value.Array()[1].String()
	return res, nil
}

func ValidateIndexConfig(index Index) error {
	if index.Dim < 1 || index.Dim > 32768 {
		return errors.New("vector dim must be in [1, 32768]")
	}
	return nil
}
