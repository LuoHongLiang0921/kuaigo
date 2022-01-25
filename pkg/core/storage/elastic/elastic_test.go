// @Description
// @Author yixia
// @Copyright 2021 sndks.com. All rights reserved.
// @LastModify 2021/1/14 5:21 下午

package elastic

import (
	"context"
	"fmt"
	"testing"

	"git.bbobo.com/framework/tabby/pkg/util/xcast"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/olivere/elastic.v6"
)

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func getES() *Elastic {
	elasticConfigs := DefaultElasticConfig()
	elastics := elasticConfigs.Build(context.TODO())
	es := elastics.Elastics["comment"]
	return es
}

func TestElastic_CreateIndex(t *testing.T) {
	type args struct {
		ctx       context.Context
		indexName string
		indexBody string
	}
	es := getES()
	indexBody := `{"settings":{"number_of_shards":1,"number_of_replicas":0},"mappings":{"demo":{"properties":{"demo_id":{"type":"integer"},"content":{"type":"text"}}}}}`
	arg := args{
		ctx:       context.TODO(),
		indexName: "demo",
		indexBody: indexBody,
	}
	tests := []struct {
		name string
		es   *Elastic
		args args
	}{
		{name: "test", es: es, args: arg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.es.CreateIndex(tt.args.ctx, tt.args.indexName, tt.args.indexBody)
			if err != nil {
				t.Errorf("CreateIndex() error = %v", err)
				return
			}
			if got == nil || !got.Acknowledged {
				t.Errorf("CreateIndex() err = %v", err)
			}
			if got != nil {
				gotString, _ := JSON.MarshalToString(got)
				fmt.Print(gotString + "\n")
			}
		})
	}
}

func TestElastic_IsExistIndex(t *testing.T) {
	type args struct {
		ctx       context.Context
		indexName string
	}
	es := getES()
	arg := args{
		ctx:       context.TODO(),
		indexName: "demo",
	}
	tests := []struct {
		name string
		es   Elastic
		args args
	}{
		{name: "test", es: *es, args: arg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.es.IsExistIndex(tt.args.ctx, tt.args.indexName)
			if err != nil {
				t.Errorf("IsExistIndex() error = %v", err)
				return
			}
			gotString, _ := JSON.MarshalToString(got)
			fmt.Print(gotString + "\n")
		})
	}
}

func TestElastic_SaveIndexBody(t *testing.T) {
	type args struct {
		ctx        context.Context
		indexName  string
		indexType  string
		indexBody  string
		indexDocId string
	}
	es := getES()
	indexBody := `{"demo_id":1,"content":"demo1"}`
	arg := args{
		ctx:        context.TODO(),
		indexName:  "demo",
		indexType:  "demo",
		indexBody:  indexBody,
		indexDocId: "10000",
	}
	tests := []struct {
		name string
		es   *Elastic
		args args
		want *IndexResponse
	}{
		{name: "test", es: es, args: arg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.es.SaveIndexBody(tt.args.ctx, tt.args.indexName, tt.args.indexType, tt.args.indexBody, tt.args.indexDocId)
			if err != nil {
				t.Errorf("SaveIndexBody() error = %v", err)
				return
			}
			if got.Status != 0 {
				t.Errorf("SaveIndexBody() got = %v", got)
			}
			gotString, _ := JSON.MarshalToString(got)
			fmt.Print(gotString + "\n")
		})
	}
}

func TestElastic_BulkSaveIndexBody(t *testing.T) {
	type args struct {
		ctx           context.Context
		indexName     string
		indexType     string
		bulkIndexBody *[]BulkIndexBody
	}
	es := getES()
	var bulkBodys []BulkIndexBody
	for i := 0; i < 10; i++ {
		bodyData := make(map[string]interface{})
		bodyData["demo_id"] = i
		bodyData["content"] = "demo_bulk" + xcast.ToString(i)
		bodyDataString, _ := JSON.MarshalToString(bodyData)
		bulkData := BulkIndexBody{
			DocId:    xcast.ToString(i),
			BodyData: bodyDataString,
		}
		bulkBodys = append(bulkBodys, bulkData)
	}
	arg := args{
		ctx:           context.TODO(),
		indexName:     "demo",
		indexType:     "demo",
		bulkIndexBody: &bulkBodys,
	}
	tests := []struct {
		name string
		es   *Elastic
		args args
	}{
		{name: "test", es: es, args: arg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.es.BulkSaveIndexBody(tt.args.ctx, tt.args.indexName, tt.args.indexType, tt.args.bulkIndexBody)
			if err != nil {
				t.Errorf("BulkSaveIndexBody() error = %v", err)
				return
			}
			gotString, _ := JSON.MarshalToString(got)
			if got.Errors {
				t.Errorf("BulkSaveIndexBody() got = %s", gotString)
			}
			fmt.Print(gotString + "\n")
		})
	}
}

func TestElastic_GetIndexBodyByDocId(t *testing.T) {
	type args struct {
		ctx        context.Context
		indexName  string
		indexType  string
		indexDocId string
	}
	es := getES()
	arg := args{
		ctx:        context.TODO(),
		indexName:  "demo",
		indexType:  "demo",
		indexDocId: "bulk_1",
	}
	tests := []struct {
		name string
		es   *Elastic
		args args
	}{
		{name: "test", es: es, args: arg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.es.GetIndexBodyByDocId(tt.args.ctx, tt.args.indexName, tt.args.indexType, tt.args.indexDocId)
			if err != nil {
				t.Errorf("GetIndexBodyByDocId() error = %v", err)
				return
			}
			gotString, _ := JSON.MarshalToString(got)
			fmt.Print(gotString + "\n")
		})
	}
}

func TestElastic_SearchIndexBody(t *testing.T) {
	type args struct {
		ctx       context.Context
		indexName string
		indexType string
		query     elastic.Query
		sorter    []elastic.Sorter
		page      int
		limit     int
	}
	es := getES()
	boolQ := elastic.NewBoolQuery()
	var sort []elastic.Sorter
	boolQ.Filter(elastic.NewRangeQuery("createTime").Gte(0), elastic.NewRangeQuery("createTime").Lte(1615340853000))
	boolQ.Must(elastic.NewMatchQuery("content", "数据"))
	scoreSort := elastic.NewScoreSort().Desc()
	sort = append(sort, scoreSort)
	createTimeSort := elastic.NewFieldSort("createTime").Desc()
	sort = append(sort, createTimeSort)
	arg := args{
		ctx:       context.TODO(),
		indexName: "comment_104",
		indexType: "comment",
		query:     boolQ,
		sorter:    sort,
		page:      1,
		limit:     10,
	}
	tests := []struct {
		name string
		es   *Elastic
		args args
	}{
		{name: "test", es: es, args: arg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.es.SearchIndexBody(tt.args.ctx, tt.args.indexName, tt.args.indexType, tt.args.query, tt.args.sorter, tt.args.page, tt.args.limit)
			if err != nil {
				t.Errorf("SearchIndexBody() error = %v", err)
				return
			}
			gotString, _ := JSON.MarshalToString(got)
			if got.Status != 0 {
				t.Errorf("SearchIndexBody() got = %v", gotString)
			}
			fmt.Print(gotString + "\n")
		})
	}
}

func TestElastic_AggregationIndexBody(t *testing.T) {
	type args struct {
		ctx             context.Context
		indexName       string
		indexType       string
		query           elastic.Query
		aggregationName string
		aggregation     elastic.Aggregation
	}
	es := getES()
	boolQ := elastic.NewBoolQuery()
	boolQ.Filter(elastic.NewRangeQuery("createTime").Gte(0), elastic.NewRangeQuery("createTime").Lte(1615282932000))
	userTermAggr := elastic.NewTermsAggregation().Field("userId")
	createTimeMaxAggr := elastic.NewMaxAggregation().Field("createTime")
	userTermAggr.SubAggregation("createTimeMaxAggr", createTimeMaxAggr)
	arg := args{
		ctx:             context.TODO(),
		indexName:       "comment_104",
		indexType:       "comment",
		query:           boolQ,
		aggregationName: "user",
		aggregation:     userTermAggr,
	}
	tests := []struct {
		name string
		es   Elastic
		args args
	}{
		{name: "comment", es: *es, args: arg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.es.AggregationIndexBody(tt.args.ctx, tt.args.indexName, tt.args.indexType, tt.args.query, tt.args.aggregationName, tt.args.aggregation)
			if err != nil {
				t.Errorf("AggregationIndexBody() error = %v", err)
				return
			}
			if got.Status != 0 || got.Error != nil {
				t.Errorf("AggregationIndexBody() got = %v", got)
			}
			terms, b := got.Aggregations.Terms("user")
			termsString, err := JSON.MarshalToString(terms)
			fmt.Print(fmt.Sprintf("aggr:%s, success:%v \n", termsString, b))
		})
	}
}

func TestElastic_UpdateIndexBodyByDocId(t *testing.T) {
	type args struct {
		ctx        context.Context
		indexName  string
		indexType  string
		updateData map[string]interface{}
		indexDocId string
	}
	es := getES()
	updateData := make(map[string]interface{})
	// updateData["demo_id"] = 30
	updateData["content"] = "demo_bulk_update_30"
	arg := args{
		ctx:        context.TODO(),
		indexName:  "demo",
		indexType:  "demo",
		updateData: updateData,
		indexDocId: "3",
	}
	tests := []struct {
		name string
		es   *Elastic
		args args
	}{
		{name: "test", es: es, args: arg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.es.UpdateIndexBodyByDocId(tt.args.ctx, tt.args.indexName, tt.args.indexType, tt.args.updateData, tt.args.indexDocId)
			if err != nil {
				t.Errorf("UpdateIndexBodyByDocId() error = %v", err)
				return
			}
			toString, err := JSON.MarshalToString(got)
			if got.Status != 0 {
				t.Errorf("UpdateIndexBodyByDocId() got = %v", toString)
			}
			fmt.Print(toString + "\n")
		})
	}
}

func TestElastic_BulkUpdateIndexBody(t *testing.T) {
	type args struct {
		ctx       context.Context
		indexName string
		indexType string
		body      *[]BulkIndexBody
	}
	es := getES()
	bulkBody := []BulkIndexBody{
		{DocId: "10000", BodyData: `{"content":"demo100"}`},
		{DocId: "5", BodyData: `{"demo_id":50}`},
	}
	arg := args{
		ctx:       context.TODO(),
		indexName: "demo",
		indexType: "demo",
		body:      &bulkBody,
	}
	tests := []struct {
		name string
		es   *Elastic
		args args
	}{
		{name: "test", es: es, args: arg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.es.BulkUpdateIndexBody(tt.args.ctx, tt.args.indexName, tt.args.indexType, tt.args.body)
			if err != nil {
				t.Errorf("BulkUpdateIndexBody() error = %v", err)
				return
			}
			toString, err := JSON.MarshalToString(got)
			if got.Errors {
				t.Errorf("BulkUpdateIndexBody() got = %v", toString)
			}
			fmt.Print(toString + "\n")
		})
	}
}

func TestElastic_DeleteIndexBodyByDocId(t *testing.T) {
	type args struct {
		ctx        context.Context
		indexName  string
		indexType  string
		indexDocId string
	}
	es := getES()
	arg := args{
		ctx:        context.TODO(),
		indexName:  "demo",
		indexType:  "demo",
		indexDocId: "0",
	}
	tests := []struct {
		name string
		es   *Elastic
		args args
	}{
		{name: "test", es: es, args: arg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.es.DeleteIndexBodyByDocId(tt.args.ctx, tt.args.indexName, tt.args.indexType, tt.args.indexDocId)
			if err != nil {
				t.Errorf("DeleteIndexBodyByDocId() error = %v", err)
				return
			}
			toString, err := JSON.MarshalToString(got)
			if got.Status != 0 {
				t.Errorf("DeleteIndexBodyByDocId() got = %v", toString)
			}
			fmt.Print(toString + "\n")
		})
	}
}

func TestElastic_DeleteIndexBody(t *testing.T) {
	type args struct {
		ctx       context.Context
		indexName string
		indexType string
		query     elastic.Query
	}
	es := getES()
	boolQ := elastic.NewBoolQuery()
	boolQ.Filter(elastic.NewRangeQuery("demo_id").Gte(0), elastic.NewRangeQuery("demo_id").Lte(5))
	arg := args{
		ctx:       context.TODO(),
		indexName: "demo",
		indexType: "demo",
		query:     boolQ,
	}
	tests := []struct {
		name string
		es   *Elastic
		args args
	}{
		{name: "test", es: es, args: arg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.es.DeleteIndexBody(tt.args.ctx, tt.args.indexName, tt.args.indexType, tt.args.query)
			if err != nil {
				t.Errorf("DeleteIndexBody() error = %v", err)
				return
			}
			toString, err := JSON.MarshalToString(got)
			if len(got.Failures) > 0 {
				t.Errorf("DeleteIndexBody() got = %v", toString)
			}
			fmt.Print(toString + "\n")
		})
	}
}
