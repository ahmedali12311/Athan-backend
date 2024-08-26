package gis

import (
	"database/sql/driver"
	"errors"

	"github.com/goccy/go-json"
)

const (
	TypePoint      = "Point"
	TypeLineString = "LineString"
)

var (
	EmptyPointString = "{\"type\":\"Point\",\"coordinates\":[]}"
	EmptyPoint       = Point{
		Type:        TypePoint,
		Coordinates: []float64{0, 0},
	}
	EmptyLineString = "{\"type\":\"LineString\",\"coordinates\":[]}"
	EmptyLine       = LineString{
		Type:        TypeLineString,
		Coordinates: [][]float64{},
	}
)

// Geospatial Data ------------------------------------------------------------

// Valuer and scanner for json objects

func jsonValuer(v any) (*string, error) {
	if v == nil {
		return &EmptyPointString, nil
	}
	var jsonStr string
	encJSONStr, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	jsonStr = string(encJSONStr)
	if jsonStr == "{}" {
		return &EmptyPointString, nil
	}
	return &jsonStr, nil
}

func jsonScanner(receiver, value any) error {
	if value == nil {
		return json.Unmarshal([]byte(EmptyPointString), receiver)
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, receiver)
}

// Point stuct ----------------------------------------------------------------

// Point Coordinates[0] is lng, Coordinates[1] is lat.
type Point struct {
	Type        string    `db:"type"        json:"type"`
	Coordinates []float64 `db:"coordinates" json:"coordinates"`
}

func (p *Point) Value() (driver.Value, error) {
	return jsonValuer(p)
}

func (p *Point) Scan(value any) error {
	return jsonScanner(p, value)
}

// Polygon stuct --------------------------------------------------------------

type Polygon struct {
	Type        *string        `db:"type"        json:"type"`
	Coordinates *[][][]float64 `db:"coordinates" json:"coordinates"`
}

func (p *Polygon) Value() (driver.Value, error) {
	return jsonValuer(p)
}

func (p *Polygon) Scan(value any) error {
	return jsonScanner(p, value)
}

// LineString stuct -----------------------------------------------------------

type LineString struct {
	Type        string      `db:"type"        json:"type"`
	Coordinates [][]float64 `db:"coordinates" json:"coordinates"`
}

func (l *LineString) Value() (driver.Value, error) {
	return jsonValuer(l)
}

func (l *LineString) Scan(value any) error {
	return jsonScanner(l, value)
}
