package sraft

import "github.com/mitchellh/mapstructure"

// AnyConfig abstract the config of other struct.
type AnyConfig struct {
	Config map[string]string `json:"config" yaml:"config"`
}

// MergeAny will merge value from x to inner config. If the data conflicts, set the value of x.
func (gc *AnyConfig) MergeAny(x any) error {
	return mapstructure.Decode(x, gc.Config)
}

// ToAny will set value from AnyConfig to Object. But the x must a point.
func (gc *AnyConfig) ToAny(x any) error {
	return mapstructure.Decode(gc.Config, x)
}
