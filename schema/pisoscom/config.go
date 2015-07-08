package pisoscom

import (
	"encoding/json"
	"os"
	"net/url"
)

type PCParamValueRange struct {
	Start int
	End   int
}

type PCParamValues struct {
	Sale         string
	ResidentRent string
	HolidaysRent string
	RentToOwn    string
	Rent         string
	Range        PCParamValueRange
}

type PCParamOperation struct {
	PlaceHolder string
	Values      PCParamValues
}

type PCParam struct {
	Operation PCParamOperation
	Province PCParamOperation
}

type PCGenerator struct {
	Url url.URL
	Params PCParam
}

type PCConfig struct {
	Generator PCGenerator
	Base      url.URL
}

func ReadConfig(fn string) (PCConfig, error) {
	file, _ := os.Open(fn)
	decoder := json.NewDecoder(file)
	cfg := PCConfig{}
	err := decoder.Decode(&cfg)

	return cfg, err
}
