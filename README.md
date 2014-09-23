# go-gae-bigquery

a nice little package to abstract the big query api. currently supports only inserting rows (queries coming soon, feel free to fork and add stuff!)

## example

the example can be found at examples/example.go. the part you want to look at is the Track function:
```go
func Track(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// create instance of big query client
	if client, err := NewClient(&c); err != nil {
		c.Errorf(err.Error())
	} else {
		// get some data to write
		rowData := GetRowData(r)

		// write the data to the table in the dataset in a specific project
		client.InsertRow(*projectID, *datasetID, *tableID, rowData)
	}
}

```