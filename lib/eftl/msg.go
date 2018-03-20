/*
 * Copyright (c) 2001-2017 TIBCO Software Inc.
 * All Rights Reserved. Confidential & Proprietary.
 * For more information, please contact:
 * TIBCO Software Inc., Palo Alto, California, USA
 *
 * $Id: msg.go 92311 2017-03-14 20:52:25Z $
 */

package eftl

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Message represents application messages that map field names to values.
type Message map[string]interface{}

// MarshalJSON encodes the message into JSON.
func (msg Message) MarshalJSON() ([]byte, error) {
	m, err := msg.encode()
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

// UnmarshalJSON decodes the message from JSON.
func (msg Message) UnmarshalJSON(b []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	msg.decode(m)
	return nil
}

func (msg Message) encode() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for k, v := range msg {
		switch v := v.(type) {
		default:
			return nil, fmt.Errorf("unsupported type for field '%s'", k)
		case string:
			m[k] = v
		case []string:
			m[k] = v
		case int, int8, int16, int32, int64:
			m[k] = v
		case []int, []int8, []int16, []int32, []int64:
			m[k] = v
		case uint, uint8, uint16, uint32, uint64:
			m[k] = v
		case []uint, []uint16, []uint32, []uint64:
			m[k] = v
		case []byte:
			m[k] = map[string]string{"_o_": base64.StdEncoding.EncodeToString(v)}
		case float32:
			m[k] = map[string]float32{"_d_": v}
		case []float32:
			s := make([]map[string]float32, 0, len(v))
			for _, t := range v {
				s = append(s, map[string]float32{"_d_": t})
			}
			m[k] = s
		case float64:
			m[k] = map[string]float64{"_d_": v}
		case []float64:
			s := make([]map[string]float64, 0, len(v))
			for _, t := range v {
				s = append(s, map[string]float64{"_d_": t})
			}
			m[k] = s
		case time.Time:
			m[k] = map[string]int64{"_m_": v.UnixNano() / 1000000}
		case []time.Time:
			s := make([]map[string]int64, 0, len(v))
			for _, t := range v {
				s = append(s, map[string]int64{"_m_": t.UnixNano() / 1000000})
			}
			m[k] = s
		case Message:
			enc, err := v.encode()
			if err != nil {
				return nil, err
			}
			m[k] = enc
		case []Message:
			s := make([]map[string]interface{}, 0, len(v))
			for _, t := range v {
				enc, err := t.encode()
				if err != nil {
					return nil, err
				}
				s = append(s, enc)
			}
			m[k] = s
		}
	}
	return m, nil
}

func (msg Message) decode(m map[string]interface{}) Message {
	for k, v := range m {
		switch v := v.(type) {
		default:
			msg[k] = v
		case float64:
			msg[k] = int64(v)
		case []interface{}:
			if len(v) > 0 {
				switch t := v[0].(type) {
				case float64:
					s := make([]int64, 0, len(v))
					for _, elem := range v {
						i, _ := elem.(float64)
						s = append(s, int64(i))
					}
					msg[k] = s
				case string:
					s := make([]string, 0, len(v))
					for _, elem := range v {
						i, _ := elem.(string)
						s = append(s, i)
					}
					msg[k] = s
				case map[string]interface{}:
					if _, exists := t["_o_"].(string); exists {
						s := make([][]byte, 0, len(v))
						for _, elem := range v {
							if el, ok := elem.(map[string]interface{}); ok {
								i, _ := el["_o_"].(string)
								d, _ := base64.StdEncoding.DecodeString(i)
								s = append(s, d)
							}
						}
						msg[k] = s
					} else if _, exists := t["_d_"].(float64); exists {
						s := make([]float64, 0, len(v))
						for _, elem := range v {
							if el, ok := elem.(map[string]interface{}); ok {
								d, _ := el["_d_"].(float64)
								s = append(s, d)
							}
						}
						msg[k] = s
					} else if _, exists := t["_m_"].(float64); exists {
						s := make([]time.Time, 0, len(v))
						for _, elem := range v {
							if el, ok := elem.(map[string]interface{}); ok {
								d, _ := el["_m_"].(float64)
								s = append(s, time.Unix(0, int64(d)*int64(time.Millisecond)))
							}
						}
						msg[k] = s
					} else {
						s := make([]Message, 0, len(v))
						for _, elem := range v {
							if el, ok := elem.(map[string]interface{}); ok {
								s = append(s, make(Message).decode(el))
							}
						}
						msg[k] = s
					}
				}
			}
		case map[string]interface{}:
			if t, exists := v["_o_"].(string); exists {
				msg[k], _ = base64.StdEncoding.DecodeString(t)
			} else if t, exists := v["_d_"].(float64); exists {
				msg[k] = t
			} else if t, exists := v["_d_"].(string); exists {
				msg[k], _ = strconv.ParseFloat(t, 64)
			} else if t, exists := v["_m_"].(float64); exists {
				msg[k] = time.Unix(0, int64(t)*int64(time.Millisecond))
			} else {
				msg[k] = make(Message).decode(v)
			}
		}
	}
	return msg
}
