package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"io/ioutil"

	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
)

var (
	repo = flag.String("repo", "", "github repository")
)

const (
	namespace       = "github"
	githubAPIPrefix = "https://api.github.com/repos/"
)

var (
	forks = prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "forks"),
		"repository forks",
		nil, nil)
	stars = prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "stars"),
		"repository stars",
		nil, nil)
)

// Exporter is containing github repository info
type Exporter struct {
	githubPage string
	repoName   string
}

// NewExporter return a new exporter
func NewExporter(githubPage, repoName string) (*Exporter, error) {
	return &Exporter{
		githubPage: githubPage,
		repoName:   repoName,
	}, nil
}

// Describe describe all the metrics ever exported
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- forks
	ch <- stars
}

// Collect fethes metrics from github api
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	rawData, err := sendRequest(e.repoName)
	if err != nil {

	}
	fmt.Println(string(rawData))
	resp := struct {
		Forks int `json:"forks"`
		Stars int `json:"stargazers_count"`
	}{}
	json.Unmarshal(rawData, &resp)
	ch <- prometheus.MustNewConstMetric(forks, prometheus.GaugeValue, float64(resp.Forks))
	ch <- prometheus.MustNewConstMetric(stars, prometheus.GaugeValue, float64(resp.Stars))
}

func (e *Exporter) getAPI(repoName string) string {
	return githubAPIPrefix + e.repoName
}

func sendRequest(repoName string) ([]byte, error) {
	url := githubAPIPrefix + repoName
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func init() {
	prometheus.MustRegister(version.NewCollector("consul_exporter"))
}

func main() {
	fmt.Println("statring Exporter...")

	flag.Parse()
	if *repo == "" {
		fmt.Fprintln(os.Stdout, "No repositories specified, exiting")
		os.Exit(0)
	}
	exporter, err := NewExporter("", *repo)
	if err != nil {
		log.Fatal("NewExporter", err)
	}
	prometheus.MustRegister(exporter)
	http.Handle("/metrics", prometheus.Handler())
	err = http.ListenAndServe(":9108", nil)
	if err != nil {
		log.Fatal("http.ListenAndServe:", err)
	}
}
