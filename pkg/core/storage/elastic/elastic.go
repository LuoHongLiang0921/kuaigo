// @Description es 常用操作

package elastic

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/ktime"
	"net"
	"net/http"
	"time"

	"gopkg.in/olivere/elastic.v6"
)

// DefaultElasticConfig 默认配置生成
// 	@Description 默认 es 配置生成
// 	@Return *Config 设置后的默认es配置
func DefaultElasticConfig() *ElasticConfigs {
	comment := ElasticConfig{
		Addrs:               []string{"http://192.168.205.237:9200"},
		MaxIdleConns:        30,
		MaxIdleConnsPerHost: 30,
		MaxConnsPerHost:     600,
		IdleConnTimeout:     ktime.Duration("90s"),
		Timeout:             ktime.Duration("60s"),
		KeepAlive:           ktime.Duration("60s"),
		Debug:               false,
		HealthCheckEnabled:  false,
		SnifferEnabled:      false,
	}
	return &ElasticConfigs{ElasticConfigs: map[string]ElasticConfig{
		"comment": comment,
	}}
}

// RawElasticConfig apollo 获取配置
// 	@Description 获取载入配置
//	@Param ctx 上下文
//	@Param key 配置key
// 	@Return *ElasticConfigs 实例后的 es 配置
func RawElasticConfig(ctx context.Context, key string) *ElasticConfigs {
	var config = DefaultElasticConfig()
	if err := conf.UnmarshalKey(key, &config.ElasticConfigs); err != nil {
		klog.KuaigoLogger.WithContext(ctx).Panic(fmt.Sprintf("unmarshal elasticConfig fail, key:%v, error:%v", key, err))
	}
	return config
}

// Build 创建 es client
// 	@Description 创建 es 多实例 client
// 	@Receiver c ElasticConfigs
//	@Param ctx 上下文
// 	@Return *Elastics 实例后的多实例 es client
func (c *ElasticConfigs) Build(ctx context.Context) *Elastics {
	if len(c.ElasticConfigs) == 0 {
		klog.KuaigoLogger.WithContext(ctx).Panic("es config empty")
	}
	elastics := make(map[string]*Elastic)
	for esName, esConfig := range c.ElasticConfigs {
		esHttpClient := http.Client{}
		esHttpClient.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   esConfig.Timeout,
				KeepAlive: esConfig.KeepAlive,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          esConfig.MaxIdleConns,
			MaxIdleConnsPerHost:   esConfig.MaxIdleConnsPerHost,
			MaxConnsPerHost:       esConfig.MaxConnsPerHost,
			IdleConnTimeout:       esConfig.IdleConnTimeout,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		if len(esConfig.Addrs) == 0 {
			klog.KuaigoLogger.WithContext(ctx).Panic(fmt.Sprintf("es config addrs empty, name:%v, config:%v", esName, esConfig))
		}
		if esConfig.Debug {
			TraceLogger.Close = false
		}
		client, err := EsClient(elastic.SetHttpClient(&esHttpClient),
			elastic.SetHealthcheck(esConfig.HealthCheckEnabled),
			elastic.SetSniff(esConfig.SnifferEnabled),
			elastic.SetURL(esConfig.Addrs...),
			elastic.SetTraceLog(TraceLogger),
		)
		if err != nil {
			klog.KuaigoLogger.WithContext(ctx).Panic(fmt.Sprintf("es client bulid fail, name:%v, config:%v, err:%v", esName, esConfig, err))
		}
		es := Elastic{
			Config: &esConfig,
			Client: client,
		}
		elastics[esName] = &es
	}
	return &Elastics{Elastics: elastics}
}

// GetElasticClient  获取 es client
// 	@Description 获取 es client
// 	@Receiver es Elastics
//	@Param ctx 上下文
//	@Param clientName es 实例别名
// 	@Return *Elastic 实例别名对应的es client 实例
// 	@Return error 错误
func (es *Elastics) GetElasticClient(ctx context.Context, clientName string) (*Elastic, error) {
	client, ok := es.Elastics[clientName]
	if !ok {
		return nil, fmt.Errorf("%v elastic client no exist", clientName)
	}
	return client, nil
}

// CloseElasticClient 关闭 es client 实例，释放资源
// 	@Description 关闭 es client 实例，释放资源
// 	@Receiver es Elastics
//	@Param ctx 上下文
//	@Param clientName es 实例别名
// 	@Return error 错误
func (es *Elastics) CloseElasticClient(ctx context.Context, clientName string) error {
	elastic, ok := es.Elastics[clientName]
	if !ok {
		return fmt.Errorf("%v elastic client no exist", clientName)
	}
	elastic.Client.Stop()
	return nil
}

// IsExistIndex
// 	@Description 验证 index 是否存在， 存在返回 index 结构
// 	@Receiver e Elastic
//	@Param ctx 上下文
//	@Param indexName index 名字
// 	@Return *EsIndex index
// 	@Return error 错误
func (e *Elastic) IsExistIndex(ctx context.Context, indexName string) (*EsIndex, error) {
	var esIndex EsIndex
	isExist, err := e.Client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return nil, err
	}
	esIndex.IsExist = isExist
	esIndex.IndexName = indexName
	if isExist {
		indexSetting, err := e.Client.IndexGetSettings(indexName).Do(ctx)
		if err != nil {
			return nil, err
		}
		esIndex.Settings = indexSetting[indexName].Settings
		indexMapping, err := e.Client.GetMapping().Index(indexName).Do(ctx)
		if err != nil {
			return nil, err
		}
		esIndex.Mappings = indexMapping[indexName].(map[string]interface{})["mappings"]
	}
	return &esIndex, nil
}

// CreateIndex 创建 es index
// 	@Description 创建 es index
//   调用示例
//   indexBody := `{"settings":{"number_of_shards":1,"number_of_replicas":1},"mappings":{"demo":{"properties":{"id":{"type":"integer"},"content":{"type":"text"}}}}}`
//   e.Client.CreateIndex("comment", "demo", indexBody)
// 	@Receiver e Elastic
//	@Param ctx 上下文
//	@Param indexName  index 名字
//	@Param indexBody  index json 结构体
// 	@Return *IndicesCreateResult 创建index应答
// 	@Return error
func (e *Elastic) CreateIndex(ctx context.Context, indexName string, indexBody string) (*IndicesCreateResult, error) {
	createRes, err := e.Client.CreateIndex(indexName).Body(indexBody).Do(ctx)
	if err != nil {
		return nil, err
	}
	return createRes, nil
}

// SaveIndexBody 保存 index 数据
// 	@Description 保存 index 数据
// 	@Receiver e Elastic
//	@Param ctx 上下文
//	@Param indexName  index 名字
//	@Param indexType  index type m名字
//	@Param indexBody  保存 index 数据json结构体
//	@Param indexDocId index 数据id，不指定时自动生成
// 	@Return *IndexResponse 保存 index 数据请求应答
// 	@Return error
func (e *Elastic) SaveIndexBody(ctx context.Context, indexName string, indexType string, indexBody string, indexDocId string) (*IndexResponse, error) {
	index := e.Client.Index().Index(indexName).Type(indexType)
	if indexDocId != "" {
		index = index.Id(indexDocId)
	}
	result, err := index.BodyJson(indexBody).Do(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// BulkSaveIndexBody 批量保存 index 数据
// 	@Description 批量保存 index 数据
// 	@Receiver e Elastic
//	@Param ctx 上下文
//	@Param indexName index 名字
//	@Param indexType index type 名字
//	@Param bulkIndexBody 批量保存 index 数据结构体
// 	@Return *BulkResponse 批量保存 index 数据应答
// 	@Return error
func (e *Elastic) BulkSaveIndexBody(ctx context.Context, indexName string, indexType string, bulkIndexBody *[]BulkIndexBody) (*BulkResponse, error) {
	bulkReq := e.Client.Bulk()
	for _, indexBody := range *bulkIndexBody {
		req := elastic.NewBulkIndexRequest().Index(indexName).Type(indexType).Doc(indexBody.BodyData)
		if indexBody.DocId != "" {
			req = req.Id(indexBody.DocId)
		}
		bulkReq = bulkReq.Add(req)
	}
	result, err := bulkReq.Do(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateIndexBodyByDocId 根据 docId 更新 index 数据
// 	@Description 根据 docId 更新 index 数据
// 	@Receiver e Elastic
//	@Param ctx 上下文
//	@Param indexName index 名字
//	@Param indexType index type 名字
//	@Param updateData index 修改数据
//	@Param indexDocId index 修改数据 id
// 	@Return *UpdateResponse 修改 index 请求应答
// 	@Return error
func (e *Elastic) UpdateIndexBodyByDocId(ctx context.Context, indexName string, indexType string, updateData map[string]interface{}, indexDocId string) (*UpdateResponse, error) {
	result, err := e.Client.Update().Index(indexName).Type(indexType).Id(indexDocId).Doc(updateData).Do(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// BulkUpdateIndexBody 批量修改 es index body
// 	@Description 批量修改 es index body
// 	@Receiver e Elastic
//	@Param ctx 上下文
//	@Param indexName index 名字
//	@Param indexType index type 名字
//	@Param bulkIndexBody 批量修改 index 数据体
// 	@Return *BulkResponse 批量修改 index 应答
// 	@Return error
func (e *Elastic) BulkUpdateIndexBody(ctx context.Context, indexName string, indexType string, bulkIndexBody *[]BulkIndexBody) (*BulkResponse, error) {
	bulkReq := e.Client.Bulk()
	for _, indexBody := range *bulkIndexBody {
		if indexBody.DocId == "" {
			continue
		}
		req := elastic.NewBulkUpdateRequest().Index(indexName).Type(indexType).Id(indexBody.DocId).ReturnSource(true).Doc(indexBody.BodyData)
		bulkReq = bulkReq.Add(req)
	}
	result, err := bulkReq.Do(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteIndexBodyByDocId 根据 docId 删除 index 数据
// 	@Description 根据 docId 删除 index 数据
// 	@Receiver e Elastic
//	@Param ctx 上下文
//	@Param indexName index 名字
//	@Param indexType index type 名字
//	@Param indexDocId index 数据 id
// 	@Return *DeleteResponse
// 	@Return error
func (e *Elastic) DeleteIndexBodyByDocId(ctx context.Context, indexName string, indexType string, indexDocId string) (*DeleteResponse, error) {
	result, err := e.Client.Delete().Index(indexName).Type(indexType).Id(indexDocId).Do(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteIndexBody 根据 query 搜索删除 index 数据
// 	@Description 根据 query 搜索删除 index 数据
// 	@Receiver e Elastic
//	@Param ctx 上下文
//	@Param indexName index 名字
//	@Param indexType index type 名字
//	@Param query 删除 index 数据的查询条件
// 	@Return *BulkIndexByScrollResponse
// 	@Return error
func (e *Elastic) DeleteIndexBody(ctx context.Context, indexName string, indexType string, query elastic.Query) (*BulkIndexByScrollResponse, error) {
	result, err := e.Client.DeleteByQuery().Index(indexName).Type(indexType).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetIndexBodyByDocId 根据 docId 获取 es 数据
// 	@Description 根据 docId 获取 es 数据
// 	@Receiver e Elastic
//	@Param ctx 上下文
//	@Param indexName index 名字
//	@Param indexType index type 名字
//	@Param indexDocId index 数据 id
// 	@Return *GetResult
// 	@Return error
func (e *Elastic) GetIndexBodyByDocId(ctx context.Context, indexName string, indexType string, indexDocId string) (*GetResult, error) {
	result, err := e.Client.Get().Index(indexName).Type(indexType).Id(indexDocId).Do(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SearchIndexBody  搜索 index 数据
// 	@Description 搜索 index 数据
//   query构建示例
//   boolQ := elastic.NewBoolQuery()
//   boolQ.Must(elastic.NewMatchQuery("content","评论"))
//   boolQ.Filter(elastic.NewRangeQuery("createTime").Gt("0"))
//   sorter构建示例
//   var sort []elastic.Sorter
//   scoreSort := elastic.NewScoreSort().Desc()
//   sort = append(sort, scoreSort)
//   createTimeSort := elastic.NewFieldSort("createTime").Desc()
//   sort = append(sort, createTimeSort)
// 	@Receiver e Elastic
//	@Param ctx 上下文
//	@Param indexName index 名字
//	@Param indexType index type 名字
//	@Param query 搜索条件
//	@Param sorter 排序字段
//	@Param page 页码
//	@Param limit 每页多少条
// 	@Return *SearchResult 搜索数据返回
// 	@Return error
func (e *Elastic) SearchIndexBody(ctx context.Context, indexName string, indexType string, query elastic.Query, sorter []elastic.Sorter, page int, limit int) (*SearchResult, error) {
	from := 0
	if page > 1 {
		from = (page - 1) * limit
	}
	result, err := e.Client.Search().Index(indexName).Type(indexType).Query(query).SortBy(sorter...).From(from).Size(limit).Do(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AggregationIndexBody 聚合查询
// 	@Description 聚合查询
//   query 构建示例
//   boolQ := elastic.NewBoolQuery()
//   boolQ.Must(elastic.NewMatchQuery("content","评论"))
//   boolQ.Filter(elastic.NewRangeQuery("createTime").Gt("0"))
//   aggregation 构建示例
//   userTermAggr := elastic.NewTermsAggregation().Field("userId")
//   createTimeMaxAggr := elastic.NewMaxAggregation().Field("createTime")
//   userTermAggr.SubAggregation("createTimeMaxAggr", createTimeMaxAggr)
//   e.Client.AggregationIndexBody("comment", "comment_104", "comment", nil, "userTermAggr", userTermAggr)
// 	@Receiver e Elastic
//	@Param ctx 上下文
//	@Param indexName index 名字
//	@Param indexType index type 名字
//	@Param query 查询条件
//	@Param aggregationName 聚合名字
//	@Param aggregation 聚合字段、类型
// 	@Return *SearchResult 聚合数据返回结果
// 	@Return error
func (e *Elastic) AggregationIndexBody(ctx context.Context, indexName string, indexType string, query elastic.Query, aggregationName string, aggregation elastic.Aggregation) (*SearchResult, error) {
	service := e.Client.Search().Index(indexName).Type(indexType)
	if query != nil {
		service = service.Query(query)
	}
	result, err := service.Aggregation(aggregationName, aggregation).Size(0).Do(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *Elastic) UpdateByQuery(ctx context.Context, indexName string, fieldName string, fieldValue string, script string, params map[string]interface{}) error {
	_, err := e.Client.UpdateByQuery().Index(indexName).
		Query(NewTermQuery(fieldName, fieldValue)).
		Script(elastic.NewScriptInline(script).Params(params)).Do(ctx)
	return err
}
func (e *Elastic) UpdateByDocId(ctx context.Context, indexName string, indexType string, indexDocId string, script string, params map[string]interface{}) error {
	_, err := e.Client.Update().Index(indexName).Type(indexType).Id(indexDocId).Script(elastic.NewScriptInline(script).Params(params)).Do(ctx)
	return err
}
