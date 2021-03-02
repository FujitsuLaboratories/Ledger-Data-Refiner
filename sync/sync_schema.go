/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"encoding/json"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/model"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/sync/analysis"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/utils"
	"github.com/pkg/errors"
	"strconv"
)

const errSchema = -1

func (sync *SyncService) schemaAnalysis() {
	logger.Info("start schema analysis")
	sync.syncSchemaStatus = true

	if err := sync.initialSchemas(); err != nil {
		logger.Error(err)
		sync.syncSchemaStatus = false
		return
	}
	if err := sync.analyzeSchema(); err != nil {
		logger.Error(err)
		sync.syncSchemaStatus = false
		return
	}

	sync.syncSchemaStatus = false
	logger.Info("finish schema analysis")
}

func (sync *SyncService) initialSchemas() error {
	logger.Debug("initial schemas from db")
	schemas, err := model.GetSchemas()
	if err != nil {
		return err
	}
	sync.storage = make(map[string][]SchemaStorage, len(schemas))
	logger.Debugf("there are %d schemas", len(schemas))
	if len(schemas) > 0 {
		for _, schema := range schemas {
			key := strconv.Itoa(schema.ChannelId) + "-" + schema.ChaincodeName
			sync.storage[key] = append(sync.storage[key], SchemaStorage{
				schemaId:    schema.ID,
				schemaArray: schema.SchemaArray,
			})
		}
	}
	return nil
}

func (sync *SyncService) reloadSchemas() error {
	logger.Debug("reload schemas from db")
	return sync.initialSchemas()
}

func (sync *SyncService) analyzeSchema() error {
	logger.Debug("get schema undetermined documents from db")
	undeterminedDocs, err := model.GetSchemaUndeterminedDocs()
	logger.Debugf("get %d schema undetermined documents from db", len(undeterminedDocs))
	if err != nil {
		return err
	}
	if len(undeterminedDocs) > 0 {
		for _, doc := range undeterminedDocs {
			logger.WithField("document_id", doc.ID).Debug("start processing document")
			if err := sync.determineSchema(doc); err != nil {
				return err
			}
			logger.WithField("document_id", doc.ID).Debug("finish processing document")
		}
		err = sync.analyzeSchema()
		if err != nil {
			return err
		}
	}

	return nil
}

func (sync *SyncService) determineSchema(doc model.Document) error {
	if sync.storage == nil {
		return errors.New("the storage is nil")
	}
	var err error
	session := model.BeginSession()

	// check if the content of document is json
	if !utils.IsJson(string(doc.Content)) {
		logger.WithField("doc_id", doc.ID).Info("the object is not a JSON schema")
		if err = model.UpdateDocSchema(session, errSchema, doc.ID); err != nil {
			session.Rollback()
			return err
		}
		// the content of document is not json, return
		session.Commit()
		return nil
	}

	docSchemaInArray, err := analysis.GetSchemaInArray(doc.Content)
	if err != nil {
		return err
	}

	compareResult := analysis.Difference
	schemas := make([]SchemaStorage, 0)
	var ok bool

	if schemas, ok = sync.storage[strconv.Itoa(doc.ChannelId)+"-"+doc.ChaincodeName]; ok {
		var schemaId int
		for _, schema := range schemas {
			schemaArray := &analysis.SchemaArray{}
			_ = json.Unmarshal([]byte(schema.schemaArray), schemaArray)
			compareResult = analysis.SchemaCompare(docSchemaInArray, *schemaArray)
			if compareResult == analysis.Difference {
				continue
			} else {
				schemaId = schema.schemaId
				break
			}
		}

		// if the result is contain, we should update the schema JSON
		if compareResult == analysis.Contain {
			logger.WithField("schema_id", schemaId).Debug("update schema with new schema")
			schemaInJson, err := analysis.GetSchemaInJson(doc.Content)
			if err != nil {
				return err
			}
			arrayBytes, _ := json.Marshal(docSchemaInArray)
			err = model.UpdateSchema(session, schemaId, string(arrayBytes), schemaInJson)
			if err != nil {
				session.Rollback()
				return err
			}
		} else if compareResult != analysis.Difference {
			// only update schema id
			err = model.UpdateDocSchema(session, schemaId, doc.ID)
			if err != nil {
				session.Rollback()
				return err
			}
		}
	}
	if len(schemas) == 0 || compareResult == analysis.Difference {
		schemaInJson, err := analysis.GetSchemaInJson(doc.Content)
		if err != nil {
			return err
		}

		arrayBytes, _ := json.Marshal(docSchemaInArray)
		// insert schema
		schemaRow := &model.Schema{
			ChannelId:     doc.ChannelId,
			ChaincodeName: doc.ChaincodeName,
			SchemaArray:   string(arrayBytes),
			SchemaJSON:    schemaInJson,
		}

		logger.WithField("content", schemaRow).Debug("insert new schema into schemas")
		err = model.InsertSchema(session, schemaRow)
		if err != nil {
			session.Rollback()
			return err
		}
		err = model.UpdateDocSchema(session, schemaRow.ID, doc.ID)
		if err != nil {
			session.Rollback()
			return err
		}
	}
	session.Commit()
	sync.reloadSchemas()

	return nil
}
