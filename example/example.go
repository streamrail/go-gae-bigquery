package example

import (
	"appengine"
	"flag"
	"github.com/streamrail/go-gae-bigquery"
	"net/http"
	"time"
)

var (
	projectID = flag.String("bq_project_id", "", "your bigquery projectid")
	datasetID = flag.String("bg_dataset", "", "Bigquery dataset to load to")
	tableID   = flag.String("bg_table", "", "Bigquery table to load to")
)

func init() {
	if *projectID == "" {
		return nil, nil, fmt.Errorf("No project id specified")
	}
	if *datasetID == "" {
		return nil, nil, fmt.Errorf("No project id specified")
	}
	if *tableID == "" {
		return nil, nil, fmt.Errorf("No project id specified")
	}
	http.HandleFunc("/track", TrackingHandler)
}

func TrackingHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	Track(w, r)

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusNoContent)
}

func Track(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// create instance of big query client
	if client, err := gobq.NewClient(&c); err != nil {
		c.Errorf(err.Error())
	} else {
		// get some data to write
		rowData := GetRowData(r)

		// write the data to the table in the dataset in a specific project
		client.InsertRow(*projectID, *datasetID, *tableID, rowData)
	}
}

func GetRowData(r *http.Request) map[string]interface{} {
	m := map[string]interface{}{
		"Category": `json:"category"`,
		"Action":   `json:"action"`,
		"Label":    `json:"label"`,
		"Date":     time.Now(),
	}
	return m
}
