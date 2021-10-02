package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/mitchellh/mapstructure"
	"github.com/oklog/run"
)

func interrupt(cancel <-chan struct{}) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-c:
		log.Printf("signal received sig=%v", sig)
		return nil
	case <-cancel:
		return errors.New("canceled")
	}
}

var (
	addr              = flag.String("addr", ":5000", "address to listen to")
	clickhouseAddress = flag.String(
		"clickhouse-address", "tcp://127.0.0.1:9000?database=tutorial", "connection string for mysql clickhouse")

	chDB *sql.DB
)

type KeyValue struct {
	Key   string
	Value interface{}
	Kind  reflect.Kind
}

type Record struct {
	Timestamp uint64 `json:"timestamp" mapstructure:"timestamp"`
	Service   string `json:"service" mapstructure:"service"`
	Kind      string `json:"kind" mapstructure:"kind"`
	Host      string `json:"host" mapstructure:"host"`

	RawMessage string `json:"message" mapstructure:"message"`

	// Used by syslog source since it doesn't write a raw json string into RawMessage
	Other map[string]interface{} `mapstructure:",remain"`

	StructuredMessages []KeyValue
	StringKeys         []string
	StringValues       []string
}

func processBody(bodyBytes []byte) (ret []*Record, err error) {
	// Is this a json marshallable body?

	var data []interface{}
	if err = json.Unmarshal(bodyBytes, &data); err != nil {
		// We can decide if we want to handle more formats next time.
		log.Printf("Failed to unmarshal to []interface{} err=%v", err)
		return
	}

	log.Printf("Found %v records", len(data))

	var records []*Record
	for _, d := range data {
		var record Record

		if err := mapstructure.Decode(d, &record); err != nil {
			log.Printf("failed to decode, dropping message; message=%v, err=%v", d, err)
			continue
		}

		if record.Timestamp == 0 || record.Service == "" || record.Kind == "" || record.Host == "" {
			log.Printf("mandatory fields not specified, dropping record; record=%v", record)
			continue
		}

		if record.RawMessage == "" {
			log.Printf("raw message empty, dropping record; record=%v", record)
			continue
		}

		// Parse for KV pairs
		var keyvalues = make(map[string]interface{})

		for k, v := range record.Other {
			keyvalues[k] = v
		}

		_ = json.Unmarshal([]byte(record.RawMessage), &keyvalues)

		unhandledKeys := make([]string, 0)
		for k, v := range keyvalues {
			keyType := reflect.TypeOf(v)
			kv := KeyValue{
				Key:   k,
				Value: v,
				Kind:  keyType.Kind(),
			}

			switch kv.Kind {
			case reflect.String:
				record.StringKeys = append(record.StringKeys, kv.Key)
				record.StringValues = append(record.StringValues, kv.Value.(string))
			default:
				unhandledKeys = append(unhandledKeys, fmt.Sprintf("key[%v] kind[%v]", kv.Key, kv.Kind))
			}

			if len(unhandledKeys) > 0 {
				log.Printf("unhandled keys=%v, record_service=%v", strings.Join(unhandledKeys, ","), record.Service)
			}
			record.StructuredMessages = append(record.StructuredMessages, kv)
		}

		records = append(records, &record)
		// log.Printf("Converted record %+v", record)
	}

	return records, nil
}

func insertIntoCH(records []*Record) (int, error) {
	tx, err := chDB.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare(
		`
		INSERT INTO logs_v1 (
			_timestamp, _service, _kind, _host, message, string.names, string.values) VALUES(
			?, ?, ?, ?, ?, ?, ?
		)
		`,
	)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for _, record := range records {
		_, err := stmt.Exec(
			record.Timestamp,
			record.Service,
			record.Kind,
			record.Host,
			record.RawMessage,
			record.StringKeys,
			record.StringValues,
		)
		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return len(records), nil
}

func ingestMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/logs", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		defer r.Body.Close()
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		processedRecords, err := processBody(bodyBytes)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		inserted, err := insertIntoCH(processedRecords)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		log.Printf("%v records inserted", inserted)

		if inserted != len(processedRecords) {
			// when can this happen actually?
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
	})

	return mux
}

func main() {
	flag.Parse()

	httpServer := &http.Server{
		Addr:    *addr,
		Handler: ingestMux(),
	}

	conn, err := sql.Open("clickhouse", *clickhouseAddress)
	if err != nil {
		panic(fmt.Errorf("connect to clickhouse: %w", err))
	}

	if err := conn.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		panic(fmt.Errorf("clickhouse ping: %w", err))
	}
	chDB = conn

	var g run.Group
	{
		cancel := make(chan struct{})
		g.Add(func() error {
			return interrupt(cancel)
		}, func(error) {
			close(cancel)
		})
	}
	{
		g.Add(func() error {
			if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
				return err
			}

			return nil
		}, func(e error) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_ = httpServer.Shutdown(ctx)
		})
	}

	if err := g.Run(); err != nil {
		log.Printf("g.Run err=%v", err)
		os.Exit(1)
	}
}
