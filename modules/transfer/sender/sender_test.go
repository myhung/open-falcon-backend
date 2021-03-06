package sender

import (
	"encoding/json"
	"reflect"
	"testing"

	cmodel "github.com/open-falcon/common/model"
)

func TestDemultiplex(t *testing.T) {
	const size = 10
	caseIn := []*cmodel.MetaData{}
	caseNqmIcmpOut := []*cmodel.MetaData{}
	caseNqmTcpOut := []*cmodel.MetaData{}
	caseNqmTcpconnOut := []*cmodel.MetaData{}
	caseGenOut := []*cmodel.MetaData{}

	for i := 0; i < size; i++ {
		if i%3 == 0 {
			fv := &cmodel.MetaData{
				Metric: "nqm-fping",
				Step:   int64(i),
			}
			caseIn = append(caseIn, fv)
			caseNqmIcmpOut = append(caseNqmIcmpOut, fv)
			caseNqmTcpOut = append(caseNqmTcpOut, fv)
			caseNqmTcpconnOut = append(caseNqmTcpconnOut, fv)
		} else {
			fv := &cmodel.MetaData{
				Metric: "test.metric.niean.1",
				Step:   int64(i),
			}
			caseIn = append(caseIn, fv)
			caseGenOut = append(caseGenOut, fv)
		}
	}

	nqmFpingItems, _, _, genericItems := Demultiplex(caseIn)
	for i, v := range nqmFpingItems {
		if v != caseNqmIcmpOut[i] {
			t.Error("Nqm item does not demultiplex properly", v)
		}
	}

	for i, v := range genericItems {
		if v != caseGenOut[i] {
			t.Error("Generic item does not demultiplex properly", v)
		}
	}
	t.Log("Nqm cases: ", nqmFpingItems, caseNqmIcmpOut)
	t.Log("Generic cases: ", genericItems, caseGenOut)
}

func createMetaData() *cmodel.MetaData {
	in := cmodel.MetaData{
		Metric:      "nqm-fping",
		Timestamp:   1460366463,
		Step:        60,
		Value:       0.000000,
		CounterType: "",
		Tags: map[string]string{
			"rttmin":               "18.64",
			"rttavg":               "21",
			"rttmax":               "26.56",
			"rttmdev":              "234.2",
			"rttmedian":            "21.5",
			"pkttransmit":          "13",
			"pktreceive":           "12",
			"dstpoint":             "test.endpoint.niean.2",
			"agent-id":             "1334",
			"agent-isp-id":         "12",
			"agent-province-id":    "13",
			"agent-city-id":        "14",
			"agent-name-tag-id":    "123",
			"agent-group-tag-ids":  "12-13-14",
			"target-id":            "2334",
			"target-isp-id":        "22",
			"target-province-id":   "23",
			"target-city-id":       "24",
			"target-name-tag-id":   "223",
			"target-group-tag-ids": "22-23-24",
		},
	}

	return &in
}

func TestConvert2NqmMetrics(t *testing.T) {
	in := cmodel.MetaData{
		Metric:      "nqm-fping",
		Timestamp:   1460366463,
		Step:        60,
		Value:       0.000000,
		CounterType: "",
		Tags: map[string]string{
			"rttmin":               "18.64",
			"rttavg":               "21",
			"rttmax":               "26.56",
			"rttmdev":              "234.2",
			"rttmedian":            "21.5",
			"pkttransmit":          "13",
			"pktreceive":           "12",
			"dstpoint":             "test.endpoint.niean.2",
			"agent-id":             "1334",
			"agent-isp-id":         "12",
			"agent-province-id":    "13",
			"agent-city-id":        "14",
			"agent-name-tag-id":    "123",
			"agent-group-tag-ids":  "12-13-14",
			"target-id":            "2334",
			"target-isp-id":        "22",
			"target-province-id":   "23",
			"target-city-id":       "24",
			"target-name-tag-id":   "223",
			"target-group-tag-ids": "22-23-24",
		},
	}
	out_ptr, _ := convert2NqmMetrics(&in)
	out := nqmMetrics{
		Rttmin:      18,
		Rttavg:      21,
		Rttmax:      26,
		Rttmdev:     234.2,
		Rttmedian:   21.5,
		Pkttransmit: 13,
		Pktreceive:  12,
	}

	if out != *out_ptr {
		t.Error("Expected output: ", out)
		t.Error("Real output:     ", *out_ptr)
	}

	in.Tags["rttmin"] = "qqqq"
	out_ptr_e, err := convert2NqmMetrics(&in)
	if out_ptr_e != nil {
		t.Error("Expected parsing error: ", err)
	}
}

func TestConvert2NqmEndpoint(t *testing.T) {
	in := cmodel.MetaData{
		Metric:      "nqm-fping",
		Timestamp:   1460366463,
		Step:        60,
		Value:       0.000000,
		CounterType: "",
		Tags: map[string]string{
			"rttmin":               "18.64",
			"rttavg":               "21",
			"rttmax":               "26.56",
			"rttmdev":              "234.2",
			"rttmedian":            "21.5",
			"pkttransmit":          "13",
			"pktreceive":           "12",
			"dstpoint":             "test.endpoint.niean.2",
			"agent-id":             "1334",
			"agent-isp-id":         "12",
			"agent-province-id":    "13",
			"agent-city-id":        "14",
			"agent-name-tag-id":    "123",
			"agent-group-tag-ids":  "12-13-14",
			"target-id":            "2334",
			"target-isp-id":        "22",
			"target-province-id":   "23",
			"target-city-id":       "24",
			"target-name-tag-id":   "223",
			"target-group-tag-ids": "22-23-24",
		},
	}
	out_ptr, _ := convert2NqmEndpoint(&in, "agent")
	out := nqmEndpoint{
		Id:          1334,
		IspId:       12,
		ProvinceId:  13,
		CityId:      14,
		NameTagId:   123,
		GroupTagIds: []int32{12, 13, 14},
	}

	if !reflect.DeepEqual(out, *out_ptr) {
		t.Error("Expected output: ", out)
		t.Error("Real output:     ", *out_ptr)

	}

	out_ptr, _ = convert2NqmEndpoint(&in, "target")
	out = nqmEndpoint{
		Id:          2334,
		IspId:       22,
		ProvinceId:  23,
		CityId:      24,
		NameTagId:   223,
		GroupTagIds: []int32{22, 23, 24},
	}

	if !reflect.DeepEqual(out, *out_ptr) {
		t.Error("Expected output: ", out)
		t.Error("Real output:     ", *out_ptr)
	}

	in.Tags["agent-id"] = "qqqq"
	out_ptr_e, err := convert2NqmEndpoint(&in, "agent")
	if out_ptr_e != nil {
		t.Error("Expected parsing error: ", err)
	}
}

func TestConvert2NqmPingItem(t *testing.T) {
	tests := []struct {
		input    *cmodel.MetaData
		expected *nqmPingItem
	}{
		{
			&cmodel.MetaData{
				Metric:      "nqm-fping",
				Timestamp:   1460366463,
				Step:        60,
				Value:       0.000000,
				CounterType: "",
				Tags: map[string]string{
					"rttmin":               "18.64",
					"rttavg":               "21.01",
					"rttmax":               "26.58",
					"rttmdev":              "234.2",
					"rttmedian":            "21.5",
					"pkttransmit":          "13",
					"pktreceive":           "12",
					"dstpoint":             "test.endpoint.niean.2",
					"agent-id":             "1334",
					"agent-isp-id":         "12",
					"agent-province-id":    "13",
					"agent-city-id":        "14",
					"agent-name-tag-id":    "123",
					"agent-group-tag-ids":  "12-13-14",
					"target-id":            "2334",
					"target-isp-id":        "22",
					"target-province-id":   "23",
					"target-city-id":       "24",
					"target-name-tag-id":   "223",
					"target-group-tag-ids": "22-23-24",
				},
			},
			&nqmPingItem{
				Timestamp: 1460366463,
				Agent: nqmEndpoint{
					Id:          1334,
					IspId:       12,
					ProvinceId:  13,
					CityId:      14,
					NameTagId:   123,
					GroupTagIds: []int32{12, 13, 14},
				},
				Target: nqmEndpoint{
					Id:          2334,
					IspId:       22,
					ProvinceId:  23,
					CityId:      24,
					NameTagId:   223,
					GroupTagIds: []int32{22, 23, 24},
				},
				Metrics: nqmMetrics{
					Rttmin:      18,
					Rttavg:      21.01,
					Rttmax:      26,
					Rttmdev:     234.2,
					Rttmedian:   21.5,
					Pkttransmit: 13,
					Pktreceive:  12,
				},
			},
		},
		{
			&cmodel.MetaData{
				Metric:      "nqm-tcpping",
				Timestamp:   1460366463,
				Step:        60,
				Value:       0.000000,
				CounterType: "",
				Tags: map[string]string{
					"rttmin":               "18.64",
					"rttavg":               "21.01",
					"rttmax":               "26.58",
					"rttmdev":              "234.2",
					"rttmedian":            "21.5",
					"pkttransmit":          "13",
					"pktreceive":           "12",
					"dstpoint":             "test.endpoint.niean.2",
					"agent-id":             "1334",
					"agent-isp-id":         "12",
					"agent-province-id":    "13",
					"agent-city-id":        "14",
					"agent-name-tag-id":    "123",
					"agent-group-tag-ids":  "12-13-14",
					"target-id":            "2334",
					"target-isp-id":        "22",
					"target-province-id":   "23",
					"target-city-id":       "24",
					"target-name-tag-id":   "223",
					"target-group-tag-ids": "22-23-24",
				},
			},
			&nqmPingItem{
				Timestamp: 1460366463,
				Agent: nqmEndpoint{
					Id:          1334,
					IspId:       12,
					ProvinceId:  13,
					CityId:      14,
					NameTagId:   123,
					GroupTagIds: []int32{12, 13, 14},
				},
				Target: nqmEndpoint{
					Id:          2334,
					IspId:       22,
					ProvinceId:  23,
					CityId:      24,
					NameTagId:   223,
					GroupTagIds: []int32{22, 23, 24},
				},
				Metrics: nqmMetrics{
					Rttmin:      18,
					Rttavg:      21.01,
					Rttmax:      26,
					Rttmdev:     234.2,
					Rttmedian:   21.5,
					Pkttransmit: 13,
					Pktreceive:  12,
				},
			},
		},
	}

	for _, v := range tests {
		got, _ := convert2NqmPingItem(v.input)
		if !reflect.DeepEqual(got, v.expected) {
			t.Error(got, "!=", v.expected)
		}
		t.Log(got, "==", v.expected)
	}
}

func TestConvert2NqmConnItem(t *testing.T) {
	tests := []struct {
		input    *cmodel.MetaData
		expected *nqmConnItem
	}{
		{
			&cmodel.MetaData{
				Metric:      "nqm-tcpconn",
				Timestamp:   1460366463,
				Step:        60,
				Value:       0.000000,
				CounterType: "",
				Tags: map[string]string{
					"time":                 "18.64",
					"dstpoint":             "test.endpoint.niean.2",
					"agent-id":             "1334",
					"agent-isp-id":         "12",
					"agent-province-id":    "13",
					"agent-city-id":        "14",
					"agent-name-tag-id":    "123",
					"agent-group-tag-ids":  "12-13-14",
					"target-id":            "2334",
					"target-isp-id":        "22",
					"target-province-id":   "23",
					"target-city-id":       "24",
					"target-name-tag-id":   "223",
					"target-group-tag-ids": "22-23-24",
				},
			},
			&nqmConnItem{
				Timestamp: 1460366463,
				Agent: nqmEndpoint{
					Id:          1334,
					IspId:       12,
					ProvinceId:  13,
					CityId:      14,
					NameTagId:   123,
					GroupTagIds: []int32{12, 13, 14},
				},
				Target: nqmEndpoint{
					Id:          2334,
					IspId:       22,
					ProvinceId:  23,
					CityId:      24,
					NameTagId:   223,
					GroupTagIds: []int32{22, 23, 24},
				},
				TotalTime: 18.64,
			},
		},
	}

	for _, v := range tests {
		got, _ := convert2NqmConnItem(v.input)
		if !reflect.DeepEqual(got, v.expected) {
			t.Error(got, "!=", v.expected)
		}
		t.Log(got, "==", v.expected)
	}

}

func TestJsonMarshal(t *testing.T) {
	in := createMetaData()
	out, _ := convert2NqmEndpoint(in, "agent")
	check, _ := json.Marshal(out)
	t.Log("JsonMarshal of agent: ", string(check))
	var dat map[string]int
	json.Unmarshal(check, &dat)

	expected := map[string]int{
		"name_tag_id": 123,
		"id":          1334,
		"isp_id":      12,
		"province_id": 13,
		"city_id":     14,
	}

	for k, v := range expected {
		if v != dat[k] {
			t.Error("Expected output: ", expected)
			t.Error("Real output:     ", dat)
		}
	}

	out, _ = convert2NqmEndpoint(in, "target")
	check, _ = json.Marshal(out)
	t.Log("JsonMarshal of target: ", string(check))
	json.Unmarshal(check, &dat)

	expected = map[string]int{
		"name_tag_id": 223,
		"id":          2334,
		"isp_id":      22,
		"province_id": 23,
		"city_id":     24,
	}

	for k, v := range expected {
		if v != dat[k] {
			t.Error("Expected output: ", expected)
			t.Error("Real output:     ", dat)
		}
	}

	m_out, _ := convert2NqmMetrics(in)
	check, _ = json.Marshal(m_out)
	t.Log("JsonMarshal of metrics: ", string(check))
	var int_dat map[string]int32
	json.Unmarshal(check, &int_dat)

	var min int32 = 18
	var max int32 = 26

	if v, p := int_dat["min"]; p {
		if v != min {
			t.Error("Expected output: ", min)
			t.Error("Real output:     ", v)
		}
	}
	if v, p := int_dat["max"]; p {
		if v != max {
			t.Error("Expected output: ", max)
			t.Error("Real output:     ", v)
		}
	}
}
