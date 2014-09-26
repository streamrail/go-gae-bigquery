package example

import (
	"appengine"
	"flag"
	"fmt"
	"github.com/streamrail/go-gae-bigquery"
	"net/http"
	"time"
)

var (
	projectID = flag.String("bq_project_id", "", "your bigquery projectid")
	datasetID = flag.String("bg_dataset", "", "Bigquery dataset to load to")
	tableID   = flag.String("bg_table", "", "Bigquery table to load to")
	buff      = gobq.NewBufferedWrite(MAX_BUFFERED)
)

const (
	MAX_BUFFERED = 3
)

func init() {
	if *projectID == "" {
		fmt.Println("No project id specified")
	}
	if *datasetID == "" {
		fmt.Println("No project id specified")
	}
	if *tableID == "" {
		fmt.Println("No project id specified")
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
		// append the row to the buffer
		if err := buff.Append(rowData); err != nil {
			c.Errorf(err.Error())
		}
		c.Infof("buffered rows: %d\n", buff.Length())
		// if the buffer is full, flush it into big query.
		// the flushing resets the buffer and you can accumulate rows again
		if buff.IsFull() {
			if err := client.InsertRows(*projectID, *datasetID, *tableID, buff.Flush()); err != nil {
				c.Errorf(err.Error())
			} else {
				c.Infof("inserted rows: %d", buff.Length())
			}
		}
	}
}

func GetRowData(r *http.Request) gobq.Row {
	m := gobq.Row{
		"Category": `json:"category"`,
		"Action":   `json:"action"`,
		"Label":    `json:"label"`,
		"Date":     time.Now(),
	}
	return m
}
