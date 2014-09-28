# go-gae-bigquery [![Build Status](https://travis-ci.org/StreamRail/go-gae-bigquery.svg?branch=master)](https://travis-ci.org/StreamRail/go-gae-bigquery)

A nice little package to abstract usage of the BigQuery service on GAE. Currently supports only inserting rows (queries coming soon, feel free to fork and add stuff!)

## usage

Import the package:

```go
import (
	"github.com/streamrail/go-gae-bigquery"
)

```
and go get it using the goapp gae command:

```bash
goapp get "github.com/streamrail/go-gae-bigquery"
```

The package is now imported under the "gobq" namespace. 

## example

Running the example:

```bash
git clone https://github.com/StreamRail/go-gae-bigquery.git
cd go-gae-bigquery
cd example-batch

goapp get "github.com/streamrail/go-gae-bigquery"
goapp serve 
```

Running tests (currently only for bufferedWrite.go, to test the client.go file we need to setup a local BigQuery dev environment):
```bash
goapp test "github.com/streamrail/go-gae-bigquery"
```

The example may be found at examples-batch/example.go. The part you want to look at is the Track function:
```go
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
```

## batching 

To improve performance, you might want to batch your inserts. A request that only appends a row to a buffer takes about 10-60ms, while a request that performs an actual inserts takes about 1.3 sec! As long as you don't mind losing some rows here and there when the instance flushes the RAM memory, you can batch your inserts by utilizing the RAM of the currently running instance. 

For this purpose the package includes a thread-safe BufferedWrite implementation, which takes care of mutex over a slice of rows, and can be used to flush a batch of rows into BigQuery in a single operation. 

Be sure to set the MAX_BUFFERED to a feasible number: there are [a few limitations](https://cloud.google.com/bigquery/streaming-data-into-bigquery#quota) for batching, they suggest not to use a MAX_BUFFERED size of more than 500 etc. 
