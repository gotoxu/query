package query

import (
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"github.com/gotoxu/assert"
)

type IntAlias int

type rudeBool bool

func (id *rudeBool) UnmarshalText(text []byte) error {
	value := string(text)
	switch {
	case strings.EqualFold("Yup", value):
		*id = true
	case strings.EqualFold("Nope", value):
		*id = false
	default:
		return errors.New("value must be yup or nope")
	}
	return nil
}

type S1 struct {
	F01 int         `url:"f1"`
	F02 *int        `url:"f2"`
	F03 []int       `url:"f3"`
	F04 []*int      `url:"f4"`
	F05 *[]int      `url:"f5"`
	F06 *[]*int     `url:"f6"`
	F07 S2          `url:"f7"`
	F08 *S1         `url:"f8"`
	F09 int         `url:"-"`
	F10 []S1        `url:"f10"`
	F11 []*S1       `url:"f11"`
	F12 *[]S1       `url:"f12"`
	F13 *[]*S1      `url:"f13"`
	F14 int         `url:"f14"`
	F15 IntAlias    `url:"f15"`
	F16 []IntAlias  `url:"f16"`
	F17 S19         `url:"f17"`
	F18 rudeBool    `url:"f18"`
	F19 *rudeBool   `url:"f19"`
	F20 []rudeBool  `url:"f20"`
	F21 []*rudeBool `url:"f21"`
}

type S2 struct {
	F01 *[]*int `url:"f1"`
}

type S19 [2]byte

func (id *S19) UnmarshalText(text []byte) error {
	buf, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	if len(buf) > len(*id) {
		return errors.New("out of range")
	}
	for i := range buf {
		(*id)[i] = buf[i]
	}
	return nil
}

func TestAll(t *testing.T) {
	v := map[string][]string{
		"f1":             {"1"},
		"f2":             {"2"},
		"f3":             {"31", "32"},
		"f4":             {"41", "42"},
		"f5":             {"51", "52"},
		"f6":             {"61", "62"},
		"f7.f1":          {"71", "72"},
		"f8.f8.f7.f1":    {"81", "82"},
		"f9":             {"9"},
		"f10.0.f10.0.f6": {"101", "102"},
		"f10.0.f10.1.f6": {"103", "104"},
		"f11.0.f11.0.f6": {"111", "112"},
		"f11.0.f11.1.f6": {"113", "114"},
		"f12.0.f12.0.f6": {"121", "122"},
		"f12.0.f12.1.f6": {"123", "124"},
		"f13.0.f13.0.f6": {"131", "132"},
		"f13.0.f13.1.f6": {"133", "134"},
		"f14":            {},
		"f15":            {"151"},
		"f16":            {"161", "162"},
		"f17":            {"1a2b"},
		"f18":            {"yup"},
		"f19":            {"nope"},
		"f20":            {"nope", "yup"},
		"f21":            {"yup", "nope"},
	}
	f2 := 2
	f41, f42 := 41, 42
	f61, f62 := 61, 62
	f71, f72 := 71, 72
	f81, f82 := 81, 82
	f101, f102, f103, f104 := 101, 102, 103, 104
	f111, f112, f113, f114 := 111, 112, 113, 114
	f121, f122, f123, f124 := 121, 122, 123, 124
	f131, f132, f133, f134 := 131, 132, 133, 134
	var f151 IntAlias = 151
	var f161, f162 IntAlias = 161, 162
	var f152, f153 rudeBool = true, false
	e := S1{
		F01: 1,
		F02: &f2,
		F03: []int{31, 32},
		F04: []*int{&f41, &f42},
		F05: &[]int{51, 52},
		F06: &[]*int{&f61, &f62},
		F07: S2{
			F01: &[]*int{&f71, &f72},
		},
		F08: &S1{
			F08: &S1{
				F07: S2{
					F01: &[]*int{&f81, &f82},
				},
			},
		},
		F09: 0,
		F10: []S1{
			S1{
				F10: []S1{
					S1{F06: &[]*int{&f101, &f102}},
					S1{F06: &[]*int{&f103, &f104}},
				},
			},
		},
		F11: []*S1{
			&S1{
				F11: []*S1{
					&S1{F06: &[]*int{&f111, &f112}},
					&S1{F06: &[]*int{&f113, &f114}},
				},
			},
		},
		F12: &[]S1{
			S1{
				F12: &[]S1{
					S1{F06: &[]*int{&f121, &f122}},
					S1{F06: &[]*int{&f123, &f124}},
				},
			},
		},
		F13: &[]*S1{
			&S1{
				F13: &[]*S1{
					&S1{F06: &[]*int{&f131, &f132}},
					&S1{F06: &[]*int{&f133, &f134}},
				},
			},
		},
		F14: 0,
		F15: f151,
		F16: []IntAlias{f161, f162},
		F17: S19{0x1a, 0x2b},
		F18: f152,
		F19: &f153,
		F20: []rudeBool{f153, f152},
		F21: []*rudeBool{&f152, &f153},
	}

	s := &S1{}
	_ = NewDecoder().Decode(v, s)

	assert.DeepEqual(t, s.F01, e.F01)

	assert.NotNil(t, s.F02)
	assert.DeepEqual(t, s.F02, e.F02)

	assert.NotNil(t, s.F03)
	assert.Len(t, s.F03, 2)
	assert.DeepEqual(t, s.F03, e.F03)

	assert.NotNil(t, s.F04)
	assert.Len(t, s.F04, 2)
	assert.DeepEqual(t, s.F04, e.F04)

	assert.NotNil(t, s.F05)
	sF05, eF05 := *s.F05, *e.F05
	assert.Len(t, sF05, 2)
	assert.DeepEqual(t, sF05, eF05)

	assert.NotNil(t, s.F06)
	sF06, eF06 := *s.F06, *e.F06
	assert.Len(t, sF06, 2)
	assert.DeepEqual(t, sF06, eF06)

	assert.NotNil(t, s.F07.F01)
	sF07, eF07 := *s.F07.F01, *e.F07.F01
	assert.Len(t, sF07, 2)
	assert.DeepEqual(t, sF07, eF07)

	assert.NotNil(t, s.F08)
	assert.NotNil(t, s.F08.F08)
	assert.NotNil(t, s.F08.F08.F07.F01)
	sF08, eF08 := *s.F08.F08.F07.F01, *e.F08.F08.F07.F01
	assert.Len(t, sF08, 2)
	assert.DeepEqual(t, sF08, eF08)

	assert.DeepEqual(t, s.F09, e.F09)

	assert.NotNil(t, s.F10)
	assert.Len(t, s.F10, 1)
	assert.Len(t, s.F10[0].F10, 2)
	sF10, eF10 := *s.F10[0].F10[0].F06, *e.F10[0].F10[0].F06
	assert.NotNil(t, sF10)
	assert.Len(t, sF10, 2)
	assert.DeepEqual(t, sF10, eF10)
	sF10, eF10 = *s.F10[0].F10[1].F06, *e.F10[0].F10[1].F06
	assert.NotNil(t, sF10)
	assert.Len(t, sF10, 2)
	assert.DeepEqual(t, sF10, eF10)

	assert.NotNil(t, s.F11)
	assert.Len(t, s.F11, 1)
	assert.Len(t, s.F11[0].F11, 2)
	sF11, eF11 := *s.F11[0].F11[0].F06, *e.F11[0].F11[0].F06
	assert.NotNil(t, sF11)
	assert.Len(t, sF11, 2)
	assert.DeepEqual(t, sF11, eF11)
	sF11, eF11 = *s.F11[0].F11[1].F06, *e.F11[0].F11[1].F06
	assert.NotNil(t, sF11)
	assert.Len(t, sF11, 2)
	assert.DeepEqual(t, sF11, eF11)

	assert.NotNil(t, s.F12)
	assert.Len(t, *s.F12, 1)
	sF12, eF12 := *(s.F12), *(e.F12)
	assert.Len(t, *sF12[0].F12, 2)
	sF122, eF122 := *(*sF12[0].F12)[0].F06, *(*eF12[0].F12)[0].F06
	assert.NotNil(t, sF122)
	assert.Len(t, sF122, 2)
	assert.DeepEqual(t, sF122, eF122)
	sF122, eF122 = *(*sF12[0].F12)[1].F06, *(*eF12[0].F12)[1].F06
	assert.NotNil(t, sF122)
	assert.Len(t, sF122, 2)
	assert.DeepEqual(t, sF122, eF122)

	assert.NotNil(t, s.F13)
	assert.Len(t, *s.F13, 1)
	sF13, eF13 := *(s.F13), *(e.F13)
	assert.Len(t, *sF13[0].F13, 2)
	sF132, eF132 := *(*sF13[0].F13)[0].F06, *(*eF13[0].F13)[0].F06
	assert.NotNil(t, sF132)
	assert.Len(t, sF132, 2)
	assert.DeepEqual(t, sF132, eF132)
	sF132, eF132 = *(*sF13[0].F13)[1].F06, *(*eF13[0].F13)[1].F06
	assert.NotNil(t, sF132)
	assert.Len(t, sF132, 2)
	assert.DeepEqual(t, sF132, eF132)

	assert.DeepEqual(t, s.F14, e.F14)

	assert.DeepEqual(t, s.F15, e.F15)

	assert.NotNil(t, s.F16)
	assert.DeepEqual(t, s.F16, e.F16)

	assert.DeepEqual(t, s.F17, e.F17)

	assert.DeepEqual(t, s.F18, e.F18)

	assert.DeepEqual(t, s.F19, e.F19)

	assert.NotNil(t, s.F20)
	assert.DeepEqual(t, s.F20, e.F20)

	assert.NotNil(t, s.F21)
	assert.DeepEqual(t, s.F21, e.F21)
}

func BenchmarkAll(b *testing.B) {
	v := map[string][]string{
		"f1":             {"1"},
		"f2":             {"2"},
		"f3":             {"31", "32"},
		"f4":             {"41", "42"},
		"f5":             {"51", "52"},
		"f6":             {"61", "62"},
		"f7.f1":          {"71", "72"},
		"f8.f8.f7.f1":    {"81", "82"},
		"f9":             {"9"},
		"f10.0.f10.0.f6": {"101", "102"},
		"f10.0.f10.1.f6": {"103", "104"},
		"f11.0.f11.0.f6": {"111", "112"},
		"f11.0.f11.1.f6": {"113", "114"},
		"f12.0.f12.0.f6": {"121", "122"},
		"f12.0.f12.1.f6": {"123", "124"},
		"f13.0.f13.0.f6": {"131", "132"},
		"f13.0.f13.1.f6": {"133", "134"},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s := &S1{}
		_ = NewDecoder().Decode(v, s)
	}
}
