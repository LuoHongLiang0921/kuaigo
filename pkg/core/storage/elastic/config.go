// @Description 配置

package elastic

import (
	"time"

	"gopkg.in/olivere/elastic.v6"
)

type (
	Client                    = elastic.Client
	IndicesCreateResult       = elastic.IndicesCreateResult
	IndexResponse             = elastic.IndexResponse
	BulkResponse              = elastic.BulkResponse
	UpdateResponse            = elastic.UpdateResponse
	BulkIndexByScrollResponse = elastic.BulkIndexByScrollResponse
	DeleteResponse            = elastic.DeleteResponse
	GetResult                 = elastic.GetResult
	SearchResult              = elastic.SearchResult
	Sorter                    = elastic.Sorter
)

var (
	EsClient                    = elastic.NewClient
	NewBoolQuery                = elastic.NewBoolQuery
	NewRangeQuery               = elastic.NewRangeQuery
	NewTermQuery                = elastic.NewTermQuery
	NewMatchQuery               = elastic.NewMatchQuery
	NewNestedQuery              = elastic.NewNestedQuery
	NewGeoDistanceQuery         = elastic.NewGeoDistanceQuery
	NewGeoBoundingBoxQuery      = elastic.NewGeoBoundingBoxQuery
	NewSuggesterGeoQuery        = elastic.NewSuggesterGeoQuery
	NewScoreSort                = elastic.NewScoreSort
	NewFieldSort                = elastic.NewFieldSort
	NewGeoDistanceSort          = elastic.NewGeoDistanceSort
	NewNestedSort               = elastic.NewNestedSort
	NewScriptSort               = elastic.NewScriptSort
	NewTermsAggregation         = elastic.NewTermsAggregation
	NewMaxAggregation           = elastic.NewMaxAggregation
	NewMinAggregation           = elastic.NewMinAggregation
	NewAvgAggregation           = elastic.NewAvgAggregation
	NewStatsAggregation         = elastic.NewStatsAggregation
	NewSumAggregation           = elastic.NewSumAggregation
	NewHistogramAggregation     = elastic.NewHistogramAggregation
	NewIPRangeAggregation       = elastic.NewIPRangeAggregation
	NewMissingAggregation       = elastic.NewMissingAggregation
	NewNestedAggregation        = elastic.NewNestedAggregation
	NewRangeAggregation         = elastic.NewRangeAggregation
	NewValueCountAggregation    = elastic.NewValueCountAggregation
	NewDateHistogramAggregation = elastic.NewDateHistogramAggregation
	NewGeoBoundsAggregation     = elastic.NewGeoBoundsAggregation
	NewGeoCentroidAggregation   = elastic.NewGeoCentroidAggregation
	NewGeoDistanceAggregation   = elastic.NewGeoDistanceAggregation
	NewGeoHashGridAggregation   = elastic.NewGeoHashGridAggregation
)

// Elastics 多集群 es
type Elastics struct {
	Elastics map[string]*Elastic
}

// Elastic 结构体
type Elastic struct {
	Config *ElasticConfig
	Client *Client
}

// ElasticConfigs 多集群 es 配置
type ElasticConfigs struct {
	ElasticConfigs map[string]ElasticConfig // 不同 es 集群配置
}

// ElasticConfig 客户端配置
type ElasticConfig struct {
	Addrs               []string      `json:"addrs"`               // es 实例配置地址
	MaxIdleConns        int           `json:"maxIdleConns"`        // http client transport 最大空闲链接
	MaxIdleConnsPerHost int           `json:"maxIdleConnsPerHost"` // http client transport 每个域名最大空闲链接
	MaxConnsPerHost     int           `json:"maxConnsPerHost"`     // http client transport 每个域名最大连接数
	IdleConnTimeout     time.Duration `json:"idleConnTimeout"`     // http client transport 空闲链接存活时间
	Timeout             time.Duration `json:"timeout"`             // http client transport Dialer 超时时间
	KeepAlive           time.Duration `json:"keepAlive"`           // http client transport Dialer 存活时间
	Debug               bool          `json:"debug"`               // 开发 debug 打印 trace 日志 默认关闭
	HealthCheckEnabled  bool          `json:"healthCheckEnabled"`  // es client 是否开启健康检查
	SnifferEnabled      bool          `json:"snifferEnabled"`      // es client 是否开启嗅探
}

// BulkIndexBody 批量操作 index
type BulkIndexBody struct {
	DocId    string // DocId 使用es默认id时不用指定值
	BodyData string // BodyData
}

// EsIndex index 结构
type EsIndex struct {
	IsExist   bool        `json:"isExist"`
	IndexName string      `json:"indexName"`
	Settings  interface{} `json:"settings"`
	Mappings  interface{} `json:"mappings"`
}
