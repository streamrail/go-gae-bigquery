# go-gae-bigquery

A nice little package to abstract usage of the BigQuery service on GAE. Currently supports only inserting rows (queries coming soon, feel free to fork and add stuff!)

## usage

Import the package:

```go
import (
	"github.com/streamrail/go-gae-bigquery"
)

```
and go get it:

```bash
go get
```

The package is now imported under the "gobq" namespace. 

## example

The example can be found at examples/example.go. the part you want to look at is the Track function:
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

To improve performance, you might want to batch your inserts. As long as you don't mind to lose some rows here and there when the instance flushes the RAM memory, you can batch your inserts by utilizing the RAM of the currently running instance. Create a slice to be used as a buffer and flush it's content into BigQuery when it reaches a certain predefined limit:

```go

const (
	MAX_BUFFERED = 10
)
func Track(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if client, err := gobq.NewClient(&c); err != nil {
		c.Errorf(err.Error())
	} else {
		rowData := GetRowData(r)
		rowsBuffer = append(rowsBuffer, rowData)
		c.Infof("buffered rows: %d\n", len(rowsBuffer))
		if len(rowsBuffer) == MAX_BUFFERED {
			if err := client.InsertRows(*projectID, *datasetID, *tableID, rowsBuffer); err != nil {
				c.Errorf(err.Error())
			} else {
				c.Infof("inserted rows: %d", len(rowsBuffer))
			}
			rowsBuffer = rowsBuffer[:0]
		}
	}
}

```