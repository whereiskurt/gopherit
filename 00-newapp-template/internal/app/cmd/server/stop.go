package server

import (
	"00-newapp-template/pkg/acme"
	"00-newapp-template/pkg/client"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
	"fmt"
)

func Stop(config *config.Config, metrics *metrics.Metrics) {
	a := client.NewAdapter(config, metrics)
	a.Config.Client.EnableLogging()
	config.Log.Infof("Sending shutdown command to: %s/shutdown", config.Client.BaseURL)

	url := fmt.Sprintf("%s/shutdown", config.Client.BaseURL)

	s := acme.NewService(config.Client.BaseURL, config.Client.SecretKey, config.Client.AccessKey)
	t := acme.NewTransport(&s)

	body, status, err := t.Get(url)
	fmt.Println(fmt.Sprintf("Status from server:%d (error:%v)\n%s", status, err, body))
}
