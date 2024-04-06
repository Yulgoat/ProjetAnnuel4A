package function

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"math/rand"
	"net/http"
	"time"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	// Initialize InfluxDB
	token := "PphN1UopXSc92TYt7XI3E0ra1NQ5OS0IGRNEEKrSq4UF4JNn2fpyWFpgBQiYGo70Q3Lp8Xo1pbZAdCgwNrRVNQ=="
	url := "http://10.42.0.1:8086"
	client := influxdb2.NewClient(url, token)

	// Close client connection when handle() exits
	defer client.Close()

	// Specify organization and bucket
	org := "Mycelium"
	bucket := "Test-bucket"

	// Create a new point with a random number as value
	rand.Seed(time.Now().UnixNano())
	value := rand.Float64() * 100 // Generates a random float between 0 and 100
	tags := map[string]string{"source": "random-number"}
	fields := map[string]interface{}{"value": value}
	point := write.NewPoint("random_measurement", tags, fields, time.Now())

	// Write the point to InfluxDB
	writeAPI := client.WriteAPIBlocking(org, bucket)
	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		fmt.Printf("Error writing data to InfluxDB: %v\n", err)
		return
	}

	fmt.Printf("Random number %.2f written to InfluxDB\n", value)
	w.WriteHeader(http.StatusOK)
}
