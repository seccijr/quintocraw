package page

import (
	"time"
	"strings"
	"errors"
	"github.com/seccijr/quintocrawl/model"
	"math"
	"strconv"
	"regexp"
)

const (
	DATA_LESS_5_Y_REGEXP = "MENOS DE 5 AÑOS"
	DATA_B_5_10_Y_REGEXP = "ENTRE 5 Y 10 AÑOS"
	DATA_B_10_20_Y_REGEXP = "ENTRE 10 Y 20 AÑOS"
	DATA_B_20_30_Y_REGEXP = "ENTRE 20 Y 30 AÑOS"
	DATA_B_30_50_Y_REGEXP = "ENTRE 30 Y 50 AÑOS"
	DATA_MORE_50_Y_REGEXP = "MÁS DE 50 AÑOS"
	DATA_B_10_20_E_REGEXP = "ENTRE 10 Y 20 €"
	DATA_B_20_40_E_REGEXP = "ENTRE 20 Y 40 €"
	DATA_B_40_60_E_REGEXP = "ENTRE 40 Y 60 €"
	DATA_B_60_80_E_REGEXP = "ENTRE 60 Y 80 €"
	DATA_B_80_100_E_REGEXP = "ENTRE 80 Y 100 €"
	DATA_MORE_100_E_REGEXP = "MÁS DE 100 €"
)


func convAge(age string) (time.Time, error) {
	var t time.Time
	switch {
	case strings.Contains(age, DATA_LESS_5_Y_REGEXP):
		t = time.Now()
	case strings.Contains(age, DATA_B_5_10_Y_REGEXP):
		t = time.Now().AddDate(-5, 0, -1)
	case strings.Contains(age, DATA_B_10_20_Y_REGEXP):
		t = time.Now().AddDate(-10, 0, -1)
	case strings.Contains(age, DATA_B_20_30_Y_REGEXP):
		t = time.Now().AddDate(-20, 0, -1)
	case strings.Contains(age, DATA_B_30_50_Y_REGEXP):
		t = time.Now().AddDate(-30, 0, -1)
	case strings.Contains(age, DATA_MORE_50_Y_REGEXP):
		t = time.Now().AddDate(-50, 0, -1)
	default:
		return t, errors.New("No matching time regexp")
	}

	return t, nil
}

func convFees(fee string) (model.PriceRange, error) {
	var r model.PriceRange
	from := model.Price{Currency: "€"}
	to := model.Price{Currency: "€"}
	switch {
	case strings.Contains(fee, DATA_B_10_20_E_REGEXP):
		from.Amount = 10.0
		to.Amount = 20.0
	case strings.Contains(fee, DATA_B_20_40_E_REGEXP):
		from.Amount = 20.0
		to.Amount = 40.0
	case strings.Contains(fee, DATA_B_40_60_E_REGEXP):
		from.Amount = 40.0
		to.Amount = 60.0
	case strings.Contains(fee, DATA_B_60_80_E_REGEXP):
		from.Amount = 60.0
		to.Amount = 80.0
	case strings.Contains(fee, DATA_B_80_100_E_REGEXP):
		from.Amount = 80.0
		to.Amount = 100.0
	case strings.Contains(fee, DATA_MORE_100_E_REGEXP):
		from.Amount = 100.0
		to.Amount = math.Inf(1)
	default:
		return r, errors.New("No matching time regexp")
	}

	r = model.PriceRange{from, to}

	return r, nil
}

func convArea(areaStr string) (model.Area, error) {
	var area model.Area
	re := regexp.MustCompile(DATA_AREA_REGEXP)
	match := re.FindStringSubmatch(areaStr)
	built, err := strconv.ParseFloat(match[1], 64)
	if err == nil {
		area.Built = built
	}
	util, err := strconv.ParseFloat(match[4], 64)
	if err == nil {
		area.Use = util
	}

	return area, nil
}
