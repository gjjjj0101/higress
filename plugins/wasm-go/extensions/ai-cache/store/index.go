package store

type Index struct {
	DataType       string `required:"false" json:"dataType"`
	Dim            uint32 `required:"true" json:"dim"`
	DistanceMethod string `required:"false" json:"distanceMethod"`
	Typ            string `required:"false" json:"type"`
	Name           string `required:"true" json:"name"`
}
