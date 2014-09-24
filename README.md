# go-gae-bigquery

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

The example can be found at examples-simple/example.go. the part you want to look at is the Track function:
```go
func Track(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// create instance of big query client
	if client, err := gobq.NewClient(&c); err != nil {
		c.Errorf(err.Error())
	} else {
		// get some data to write
		rowData := GetRowData(r)

		// yoo can insert just one row, or multiple rows at the same operation
		rows := []map[string]interface{}{rowData}

		// write the data to the table in the dataset in a specific project
		client.InsertRow(*projectID, *datasetID, *tableID, rows)
	}
}
```

## batching 

To improve performance, you might want to batch your inserts. A request that only appends a row to a buffer takes about 60ms, while a request that performs an actual inserts takes about 1.3 sec! As long as you don't mind losing some rows here and there when the instance flushes the RAM memory, you can batch your inserts by utilizing the RAM of the currently running instance. For this purpose the package includes a thread-safe BufferedWrite implementation, which takes care of mutex over a slice of rows, and can be used to flush a batch of rows into BigQuery in a single operation. 


The following example flushes the buffer after 3 rows have been appended (a complete example can be found at examples-batching/example.go)

```go


const (
	MAX_BUFFERED = 3
)

var (
	buff = gobq.NewBufferedWrite(MAX_BUFFERED)
)

fu
func Track(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if client, err := gobq.NewClient(&c); err != nil {
		c.Errorf(err.Error())
	} else {
		rowData := GetRowData(r)
		if err := buff.Append(rowData); err != nil {
			c.Errorf(err.Error())
		}
		c.Infof("buffered rows: %d\n", buff.Length())
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