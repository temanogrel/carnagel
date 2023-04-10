package performer

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v5"
)

const (
	performerIndexName = "infinity-performers"
	aliasesTypeName    = "aliases"

	performerIndexMapping = `{
  "settings": {
    "number_of_shards": 5,
    "number_of_replicas": 1,
    "analysis": {
      "analyzer": {
        "lowercase": {
          "tokenizer": "stage_name_ngram",
          "token_filters": [
            "lowercase"
          ]
        }
      },
       "tokenizer": {
        "stage_name_ngram": {
          "type": "edge_ngram",
          "min_gram": 3,
          "max_gram": 20
        }
      }
    }
  },
  "mappings": {
    "aliases": {
      "properties": {
        "performerId": {
          "type": "keyword"
        },
        "stageName": {
          "type": "text",
          "analyzer": "lowercase",
          "fields": {
            "keyword": {
              "type": "keyword"
            }
          }
        }
      }
    }
  }
}`
)

func NewElasticSearchService(app *infinity.Application) infinity.PerformerSearchService {
	return &performerSearchService{
		app: app,
		log: app.Logger.WithField("component", "PerformerSearchService"),
	}
}

type performerSearchService struct {
	app *infinity.Application
	log logrus.FieldLogger

	indexRebuilding bool
}

func (service *performerSearchService) Run(ctx context.Context) {
	exists, err := service.app.ElasticSearch.IndexExists(performerIndexName).Do(context.Background())
	if err != nil {
		service.log.
			WithError(err).
			Error("Failed to check if the index exists")

		return
	}

	if !exists {
		if err := service.rebuildIndex(ctx); err == nil {
			service.log.
				WithError(err).
				Error("Failed to rebuild index")
		}

		return
	}

	if service.app.Consul.GetBool("infinity/elasticsearch-rebuild-index", false) {
		service.log.Info("Force rebuilding index")

		if _, err := service.app.ElasticSearch.DeleteIndex(performerIndexName).Do(ctx); err != nil {
			service.log.
				WithError(err).
				Error("failed to delete index")

			return
		}

		if err := service.rebuildIndex(ctx); err != nil {
			service.log.
				WithError(err).
				Error("Failed to rebuild index")
		}
	}
}

func (service *performerSearchService) rebuildIndex(ctx context.Context) error {
	log := service.log.WithField("operation", "PerformerElasticSearchRebuildIndex")
	log.Infof("Started rebuilding the index")

	service.indexRebuilding = true

	defer func() {
		service.indexRebuilding = false

		log.Info("Finished rebuilding index")
	}()

	_, err := service.app.ElasticSearch.
		CreateIndex(performerIndexName).
		BodyString(performerIndexMapping).
		Do(context.Background())

	if err != nil {
		return errors.Wrap(err, "Failed to set the index mapping")
	}

	var lastSeen time.Time
	var performers []infinity.Performer

	for {
		err := service.app.DB.
			WithContext(ctx).
			Model(&performers).
			Where("created_at > ?", lastSeen).
			Limit(5000).
			OrderExpr("created_at ASC").
			Select()

		if err != nil {
			return errors.Wrap(err, "Failed to retrieve performers")
		}

		if len(performers) == 0 {
			return nil
		}

		bulkService := service.app.ElasticSearch.Bulk()

		for _, performer := range performers {
			for _, stageName := range performer.Aliases {
				bulkService.Add(
					elastic.
						NewBulkIndexRequest().
						Index(performerIndexName).
						Type(aliasesTypeName).
						Doc(elasticPerformerDocument{PerformerId: performer.Uuid, StageName: strings.ToLower(stageName)}),
				)
			}

			lastSeen = performer.CreatedAt
		}

		if _, err = bulkService.Do(ctx); err != nil {
			return errors.Wrap(err, "Failed to bulk import performers")
		}

		// reset the performers array because go-pg will just append
		performers = make([]infinity.Performer, 0)
	}
}

func (service *performerSearchService) AddPerformerAlias(performerId uuid.UUID, stageName string) error {
	if service.indexRebuilding {
		return nil
	}

	if exists, err := service.PerformerAliasExists(performerId, stageName); err != nil || exists {
		return err
	}

	_, err := service.app.ElasticSearch.Index().
		Index(performerIndexName).
		Type(aliasesTypeName).
		BodyJson(elasticPerformerDocument{PerformerId: performerId, StageName: strings.ToLower(stageName)}).
		Do(context.Background())

	return err
}

func (service *performerSearchService) PerformerAliasExists(performerId uuid.UUID, stageName string) (bool, error) {
	query := elastic.NewBoolQuery()
	query.Filter(
		elastic.NewTermQuery("performerId", performerId),
		elastic.NewTermQuery("stageName.keyword", strings.ToLower(stageName)),
	)

	resp, err := service.app.ElasticSearch.Search().
		Index(performerIndexName).
		Type(aliasesTypeName).
		Query(query).
		Size(1).
		Do(context.Background())

	if err != nil {
		return false, errors.Wrap(err, "Failed to check if the alias already exists")
	}

	return resp.Hits.TotalHits > 0, nil
}

func (service *performerSearchService) Matching(criteria *infinity.PerformerRepositoryCriteria) ([]uuid.UUID, int, error) {
	if service.indexRebuilding {
		return []uuid.UUID{}, 0, errors.New("Index is currently building, search is unavailable")
	}

	var result []uuid.UUID

	query := elastic.NewBoolQuery()
	query.Should(
		// Prefix based queries are slightly more important
		elastic.
			NewPrefixQuery("stageName", strings.ToLower(criteria.StageName)).
			Boost(1.3),

		// This should match against the n-gram stuff
		elastic.
			NewTermQuery("stageName", strings.ToLower(criteria.StageName)),

		// exact matches
		elastic.
			NewTermQuery("stageName.keyword", strings.ToLower(criteria.StageName)).
			Boost(2.5),
	)

	resp, err := service.app.ElasticSearch.Search().
		Index(performerIndexName).
		Type(aliasesTypeName).
		Query(query).
		Size(criteria.Limit).
		From(criteria.Offset).
		Do(context.Background())

	service.log.
		WithField("stageName", strings.ToLower(criteria.StageName)).
		WithField("limit", criteria.Limit).
		WithField("offset", criteria.Offset).
		WithField("duration", resp.TookInMillis).
		Debug("Performed a stage name search")

	if err != nil {
		return result, 0, errors.Wrap(err, "Failed to check if the alias already exists")
	}

	for _, hit := range resp.Hits.Hits {
		var performerResultDoc elasticPerformerDocument

		if err := json.Unmarshal(*hit.Source, &performerResultDoc); err != nil {
			return result, 0, errors.Wrap(err, "Failed to decode json from elasticsearch")
		}

		result = append(result, performerResultDoc.PerformerId)
	}

	return result, int(resp.Hits.TotalHits), nil
}
