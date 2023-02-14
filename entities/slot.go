package entities

type Slot struct {
	Changes []Change `json:"change"`
}

type Change struct {
	Kind         string        `json:"kind"`
	Schema       string        `json:"schema"`
	Table        string        `json:"table"`
	ColumnNames  []string      `json:"columnnames"`
	ColumnValues []interface{} `json:"columnvalues"`
}

type Data struct {
	Data string
}
