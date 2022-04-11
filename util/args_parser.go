package util

import (
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/labels"
)

// ParseArgs helper to parse environment variables used to configure the package
func ParseArgs() (labels.Selector, string, int, int, error) {
	viper.SetEnvPrefix("PODCHAOSMONKEY")
	viper.BindEnv("annotations")
	viper.BindEnv("interval_seconds")
	viper.BindEnv("namespace")
	viper.BindEnv("grace_period_seconds")

	annotations, err := labels.Parse(viper.GetString("annotations"))
	if err != nil {
		return nil, "", 0, 0, err
	}

	return annotations, viper.GetString("namespace"), viper.GetInt("interval_seconds"), viper.GetInt("grace_period_seconds"), err
}
