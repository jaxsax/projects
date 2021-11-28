package queryrunner

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/scylladb/go-set/strset"
)

func sanitize(s string) string {
	replaceAllWrapper := func(ch, replacement string) func(s string) string {
		return func(s string) string {
			return strings.ReplaceAll(s, ch, replacement)
		}
	}
	fns := []func(string) string{
		strings.TrimSpace,
		replaceAllWrapper("\n", ""),
		replaceAllWrapper("\t", ""),
	}

	var tmp = s
	for _, fn := range fns {
		tmp = fn(tmp)
	}

	return tmp
}

type TableInfo struct {
	Database string
	Table    string
}

type Query struct {
	Raw          string
	From         time.Time
	To           time.Time
	Table        *TableInfo
	ReturnFields []string
}

var (
	validSysFields = strset.New(
		"service",
		"host",
		"kind",
		"message",
	)
)

func (q *Query) ListFieldStatements(serviceName string) []string {
	types := []string{
		"string",
		"float64",
	}

	var items []string
	for _, typ := range types {
		s := fmt.Sprintf("SELECT distinct(arrayJoin(flatten(groupArray(%s.names)))) FROM %s.%s WHERE _service = '%s'", typ, q.Table.Database, q.Table.Table, serviceName)
		items = append(items, s)
	}

	return items
}

func (q *Query) GetMessageStatement() (string, error) {
	rq, err := NewQuery(q.Raw)
	if err != nil {
		return "", err
	}

	conds := []string{}
	for key, value := range rq.StringFields {
		s := fmt.Sprintf("`string.values`[indexOf(`string.names`, '%s')] = '%v'", key, value)
		conds = append(conds, s)
	}

	for key, value := range rq.Float64Fields {
		s := fmt.Sprintf("`float64.values`[indexOf(`float64.names`, '%s')] = %v", key, value)
		conds = append(conds, s)
	}
	for key, value := range rq.SysFields {
		if !validSysFields.Has(key) {
			continue
		}

		if key == "message" {
			s := fmt.Sprintf("multiSearchAny(`%s`, ['%v'])", key, value)
			conds = append(conds, s)
			continue
		}

		s := fmt.Sprintf("_%s = '%v'", key, value)
		conds = append(conds, s)
	}

	condStr := ""
	if len(conds) > 0 {
		condStr = strings.Join(conds, " AND ")
	}

	return fmt.Sprintf(
		"SELECT %s FROM %s.%s WHERE _timestamp >= %d AND _timestamp <= %d AND %s",
		strings.Join(q.ReturnFields, ", "),
		q.Table.Database, q.Table.Table,
		q.From.Unix(), q.To.Unix(),
		condStr,
	), nil
}

func doClickhouseRequest(addr, query string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}

	params := url.Values{}
	params.Add("query", query)
	req.URL.RawQuery = params.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do request")
	}

	return resp.Body, nil
}

func Run() {
	clickhouseAddress := "https://clickhouse.internal.jaxsax.co"

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		rawQuery := sanitize(scanner.Text())
		if rawQuery == "" {
			continue
		}
		log.Printf("received query: '%v'", rawQuery)

		to := time.Now()
		from := to.Add(-time.Duration(10) * time.Minute)
		var query Query
		if strings.HasPrefix(rawQuery, ".fields") {
			fieldsLen := len(".fields")
			if len(rawQuery) <= fieldsLen {
				log.Printf("please specify service name, else clickhouse dies")
				continue
			}

			serviceName := rawQuery[fieldsLen+1:]
			query = Query{
				Table: &TableInfo{
					Database: "logging",
					Table:    "logs_v1",
				},
			}

			for _, stmt := range query.ListFieldStatements(serviceName) {
				log.Println(stmt)
				rdr, err := doClickhouseRequest(clickhouseAddress, stmt)
				if err != nil {
					log.Printf("list field request stmt=%v, err=%v", stmt, err)
					continue
				}
				defer rdr.Close()
				dataRaw, err := io.ReadAll(rdr)
				if err != nil {
					continue
				}

				fmt.Println(string(dataRaw))
			}
		} else {
			query = Query{
				Raw: rawQuery,
				ReturnFields: []string{
					"fromUnixTimestamp(_timestamp)", "_service", "_host", "_kind", "message",
				},
				Table: &TableInfo{
					Database: "logging",
					Table:    "logs_v1",
				},
				From: from,
				To:   to,
			}
			stmt, err := query.GetMessageStatement()
			if err != nil {
				log.Printf("query: '%v' err=%v", query.Raw, err)
				continue
			}

			log.Println(stmt)
			respReader, err := doClickhouseRequest(clickhouseAddress, stmt)
			if err != nil {
				log.Printf("do request err=%v", err)
				continue
			}
			defer respReader.Close()

			rawBytes, err := io.ReadAll(respReader)
			if err != nil {
				log.Printf("error reading body: err=%v", err)
				continue
			}

			print(string(rawBytes))
		}

	}

}
