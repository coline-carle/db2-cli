package dblayout

import (
	"encoding/xml"
	"fmt"
)

// Field ...
type Field struct {
	XMLName xml.Name `xml:"Field"`
	Name    string   `xml:"Name,attr"`
	Type    string   `xml:"Type,attr"`
}

// Table ...
type Table struct {
	XMLName xml.Name `xml:"Table"`
	Name    string   `xml:"Name,attr"`
	Build   int      `xml:"Build,attr"`
	Fields  []Field  `xml:"Field"`
}

// DBFilesClient ...
type DBFilesClient struct {
	XMLName xml.Name `xml:"DBFilesClient"`
	Tables  []Table  `xml:"Table"`
}

// LoadLayout ...
func LoadLayout(version int) (dbFilesClient *DBFilesClient, err error) {
	filename := fmt.Sprintf("data/dblayout/%d.xml", version)

	data, err := Asset(filename)
	if err != nil {
		return nil, err
	}

	dbFilesClient = &DBFilesClient{}

	err = xml.Unmarshal(data, &dbFilesClient)
	if err != nil {
		return nil, err
	}
	return dbFilesClient, nil
}
