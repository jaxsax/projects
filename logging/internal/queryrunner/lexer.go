package queryrunner

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// What should this be doing??
// We want something to convert from the raw queries ie: things that you'll write in the query box no time ranges
// status_code:502 service:envoy
// authority:clickhouse-internal.jaxsax.co
// https://lucene.apache.org/core/2_9_4/queryparsersyntax.html
// This is kinda inspired by lucene syntax
// any field-like things matching this syntax KEY:VALUE will be converted into a typed search and removed out of consideration
// anything else, will go into a multiSearchAny message
// this should support negative searches as well, so things like
// -status_code:502 will exclude any 502
// likewise, for message this will be -"hello" will exclude anything where substring contains "hello"

// v1: Most simple form requirements
// substring search
// keyword search
//  - anything where everything are numbers, search for it in the float64 cols
//  - anything else, do a text search; exact match

var (
	ErrUnbalancedKeywords = errors.New("unbalanced keywords")
)

type RawQuery struct {
	StringFields  map[string]string
	Float64Fields map[string]float64
	SysFields     map[string]string
	FreeForm      string
}

func allDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

func characters(s string) bool {
	s = strings.ToLower(s)
	for _, c := range s {
		if c > 'a' && c < 'z' {
			continue
		}
		if c == '.' {
			continue
		}
	}
	return true
}

func NewQuery(query string) (*RawQuery, error) {
	query = sanitize(query)
	log.Printf("'%v'", query)

	rq := &RawQuery{
		StringFields:  make(map[string]string),
		Float64Fields: make(map[string]float64),
		SysFields:     make(map[string]string),
	}
	splits := strings.Split(query, " ")
	for _, split := range splits {
		if strings.HasPrefix(split, "@") {
			// @service:envoy
			if len(split) <= 1 {
				continue
			}

			// Remove the '@' ch
			atValue := split[1:]

			// Split into a key value
			parts := strings.Split(atValue, ":")
			if len(parts) > 2 {
				continue
			}

			rq.SysFields[parts[0]] = parts[1]
		} else if strings.Contains(split, ":") {
			parts := strings.Split(split, ":")

			if len(parts)%2 != 0 {
				return nil, ErrUnbalancedKeywords
			}

			for i := 0; i < len(parts); i += 2 {
				k := parts[i]
				v := parts[i+1]

				if allDigits(v) {
					f, err := strconv.ParseFloat(v, 64)
					if err != nil {
						return nil, fmt.Errorf("parse float: %w", err)
					}

					rq.Float64Fields[k] = f
					continue
				}

				if characters(v) {
					rq.StringFields[k] = v
				}
			}
		}
	}

	return rq, nil
}
