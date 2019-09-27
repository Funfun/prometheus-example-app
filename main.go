package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var ()

func main() {
	bind := ""
	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagset.StringVar(&bind, "bind", ":8080", "The socket to bind to.")
	flagset.Parse(os.Args[1:])

	exporter, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		log.Fatal(err)
	}
	view.Register(
		ochttp.ServerRequestCountView,
		ochttp.ServerRequestBytesView,
		ochttp.ServerResponseBytesView,
		ochttp.ServerLatencyView,
		ochttp.ServerRequestCountByMethod,
		&view.View{
			Name:        "opencensus.io/http/server/response_count_by_status_code_and_path",
			Description: "Server response count by status code",
			TagKeys:     []tag.Key{ochttp.StatusCode, ochttp.KeyServerRoute},
			Measure:     ochttp.ServerLatency,
			Aggregation: view.Count(),
		})

	view.SetReportingPeriod(1 * time.Second)

	datetimeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(time.Now().Format("2006-01-02T15:04:05Z")))
	})
	internalErr := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	http.Handle("/", ochttp.WithRouteTag(datetimeHandler, "/"))
	http.Handle("/err", ochttp.WithRouteTag(internalErr, "/err"))

	http.Handle("/metrics", exporter)
	log.Fatal(http.ListenAndServe(bind, &ochttp.Handler{}))
}
