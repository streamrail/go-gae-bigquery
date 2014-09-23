package gobq

import (
	"appengine"
	"fmt"

	"code.google.com/p/goauth2/appengine/serviceaccount"
	"code.google.com/p/google-api-go-client/bigquery/v2"
)

type Client struct {
	Service *bigquery.Service
}

func connect(c *appengine.Context) (s *bigquery.Service, e error) {
	if client, err := serviceaccount.NewClient(*c, bigquery.BigqueryScope); err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	} else {
		if service, err := bigquery.New(client); err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		} else {
			return service, nil
		}
	}
}

func NewClient(c *appengine.Context) (*Client, error) {
	if s, err := connect(c); err != nil {
		return nil, err
	} else {
		client := &Client{
			Service: s,
		}
		return client, nil
	}
}

func (c *Client) InsertRow(projectID string, datasetID string, tableID string, rowData map[string]interface{}) error {
	jsonRows := make(map[string]bigquery.JsonValue)
	for key, value := range rowData {
		jsonRows[key] = bigquery.JsonValue(value)
	}
	rows := []*bigquery.TableDataInsertAllRequestRows{
		{
			Json: jsonRows,
		},
	}

	insertRequest := &bigquery.TableDataInsertAllRequest{Rows: rows}

	result, err := c.Service.Tabledata.InsertAll(projectID, datasetID, tableID, insertRequest).Do()
	if err != nil {
		return fmt.Errorf("Error inserting row: %v", err)
	}

	if len(result.InsertErrors) > 0 {
		errstr := fmt.Sprintf("Insertion for %d rows failed\n", len(result.InsertErrors))
		for _, errors := range result.InsertErrors {
			for _, errorproto := range errors.Errors {
				errstr += fmt.Sprintf("Error inserting row %d: %+v\n", errors.Index, errorproto)
			}
		}
		return fmt.Errorf(errstr)
	}
	return nil
}
