package flag

import (
	"testing"
	"time"

	"github.com/jbsmith7741/trial"
)

func TestNew(t *testing.T) {
	type tFlag struct {
		//Name  string
		Usage string
		Def   string
	}
	type Aint int
	type Astring string
	type AFloat64 float64
	fn := func(args ...interface{}) (interface{}, error) {
		f, err := New(args[0])
		if err != nil {
			return nil, err
		}
		result := make(map[string]*tFlag)
		for key := range f.values {
			v := f.flagSet.Lookup(key)
			if v == nil {
				result[key] = nil
				continue
			}
			result[key] = &tFlag{Usage: v.Usage, Def: v.DefValue}
		}
		return result, nil
	}
	cases := trial.Cases{
		"ints": {
			Input: &struct {
				Int   int
				Int8  int8
				Int16 int16
				Int32 int32
				Int64 int64
			}{
				Int:   1,
				Int8:  2,
				Int16: 3,
				Int32: 4,
				Int64: 5,
			},
			Expected: map[string]*tFlag{
				"int":    {Def: "1"},
				"int-8":  {Def: "2"},
				"int-16": {Def: "3"},
				"int-32": {Def: "4"},
				"int-64": {Def: "5"},
			},
		},
		"uints": {
			Input: &struct {
				Uint   uint
				Uint8  uint8
				Uint16 uint16
				Uint32 uint32
				Uint64 uint64
			}{
				Uint:   6,
				Uint8:  7,
				Uint16: 8,
				Uint32: 9,
				Uint64: 10,
			},
			Expected: map[string]*tFlag{
				"uint":    {Def: "6"},
				"uint-8":  {Def: "7"},
				"uint-16": {Def: "8"},
				"uint-32": {Def: "9"},
				"uint-64": {Def: "10"},
			},
		},
		"floats": {
			Input: &struct {
				Float32 float32
				Float64 float64
			}{
				Float32: 1.0,
				Float64: 2.2,
			},
			Expected: map[string]*tFlag{
				"float-32": {Def: "1"},
				"float-64": {Def: "2.2"},
			},
		},
		"bool": {
			Input: &struct{ Bool bool }{Bool: true},
			Expected: map[string]*tFlag{
				"bool": {Def: "true"},
			},
		},
		"string": {
			Input: &struct{ String string }{String: "Hello"},
			Expected: map[string]*tFlag{
				"string": {Def: "Hello"},
			},
		},
		"with tags": {
			Input: &struct {
				Int  int    `flag:"Count" comment:"number of people in a room"`
				Name string `flag:"-" comment:"ignore me"`
			}{
				Int:  10,
				Name: "Bob",
			},
			Expected: map[string]*tFlag{
				"Count": {Def: "10", Usage: "number of people in a room"},
			},
		},
		"time": {
			Input: &struct {
				Time     time.Time
				CTime    time.Time `flag:"ctime" fmt:"2006-01-02"`
				WaitTime time.Duration
			}{
				Time:     trial.TimeDay("2019-01-02"),
				CTime:    trial.TimeDay("2019-01-02"),
				WaitTime: time.Hour,
			},
			Expected: map[string]*tFlag{
				"time":      {Def: "2019-01-02T00:00:00Z"},
				"ctime":     {Def: "2019-01-02"},
				"wait-time": {Def: "1h0m0s"},
			},
		},
		"text marshaler": {
			Input: &struct {
				MyStruct marshalStruct
				IStruct  marshalStruct `flag:"-"`
			}{
				MyStruct: marshalStruct{value: "a"},
				IStruct:  marshalStruct{value: "b"},
			},
			Expected: map[string]*tFlag{
				"my-struct": {Def: "a"},
			},
		},
		/* "alias no marshaler": {
			Input: &struct{ Int Aint }{Int: Aint(3)},
			Expected: map[string]*tFlag{
				"int": {Def: "3"},
			},
		}, */
		"pointers": {
			Input: &struct {
				Int *int
				//	Uint     *uint
				String   *string
				MyStruct *marshalStruct
			}{
				Int:      trial.IntP(1),
				String:   trial.StringP("a"),
				MyStruct: &marshalStruct{value: "c"},
			},
			Expected: map[string]*tFlag{
				"int":       {Def: "1"},
				"string":    {Def: "a"},
				"my-struct": {Def: "c"},
			},
		},
	}
	trial.New(fn, cases).SubTest(t)
}

type marshalStruct struct {
	value string
}

func (m marshalStruct) MarshalText() ([]byte, error) {
	return []byte(m.value), nil
}

func (m *marshalStruct) UnmarshalText(b []byte) error {
	m.value = string(b)
	return nil
}
