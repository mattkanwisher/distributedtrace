package zipkin

import (
	"github.com/hyperworks/influxdb/client"
)

type influxOutput struct {
	*client.Client
	clientConfig *client.ClientConfig
	config       *Config
	seriesName   string
}

// NewInfluxOutput() creates an Output that converts ZipKin spans to InfluxDB
// (http://influxdb.com) series/data points and writes it to the specified InfluxDB
// address and database.
func NewInfluxOutput(config *Config, address, database, username, password, seriesName string) (o Output, e error) {
	defer autoRecover(&e)
	clientConfig := &client.ClientConfig{
		Host:     address,
		Username: username,
		Password: password,
		Database: database,
		IsSecure: false,
	}

	cl, e := client.New(clientConfig)
	noError(e)

	if e = cl.Ping(); e != nil {
		return nil, e
	}

	return &influxOutput{cl, clientConfig, config, seriesName}, nil
}

func (inf *influxOutput) Write(result OutputMap) error {
	// convert to influx series
	keys := make([]string, 0, len(result))
	values := make([]interface{}, 0, len(result))
	for k, v := range result {
		keys = append(keys, k)
		values = append(values, v)
	}

	series := &client.Series{
		Name:    inf.seriesName,
		Columns: keys,
		Points:  [][]interface{}{values},
	}

	return inf.WriteSeries([]*client.Series{series})
}
