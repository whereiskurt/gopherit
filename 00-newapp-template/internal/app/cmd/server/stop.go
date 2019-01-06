package server

import (
	"00-newapp-template/pkg/acme"
	"00-newapp-template/pkg/client"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
	"fmt"
)

// Stop visits the specific '/shutdown' URL beginning the clean server shutdown
func Stop(config *config.Config, metrics *metrics.Metrics) {
	a := client.NewAdapter(config, metrics)
	a.Config.Client.EnableLogging()
	config.Log.Debugf("Sending shutdown command to: %s/shutdown", config.Client.BaseURL)

	url := fmt.Sprintf("%s/shutdown", config.Client.BaseURL)

	s := acme.NewService(config.Client.BaseURL, config.Client.SecretKey, config.Client.AccessKey)
	t := acme.NewTransport(&s)

	body, status, err := t.Get(url)
	if err != nil {
		config.Log.Infof("Server at '%s' was not running or cannot be reached: error: '%v'", url, err)
		return
	}

	fmt.Println(fmt.Sprintf("Success [%d]!\n%s", status, body))
	return
}
