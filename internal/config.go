package internal

import "github.com/phramz/go-healthcheck/pkg/contract"

var _ contract.Config = (*config)(nil)

func NewConfig() contract.Config {
	return &config{
		stringValues:   make(map[string]string),
		boolValues:     make(map[string]bool),
		intValues:      make(map[string]int),
		intSliceValues: make(map[string][]int),
	}
}

type config struct {
	stringValues   map[string]string
	boolValues     map[string]bool
	intValues      map[string]int
	intSliceValues map[string][]int
}

func (c *config) clone() *config {
	return &config{
		stringValues:   c.stringValues,
		boolValues:     c.boolValues,
		intValues:      c.intValues,
		intSliceValues: c.intSliceValues,
	}
}

func (c *config) IntSlice(key string, defaultValue []int) []int {
	if _, ok := c.intSliceValues[key]; ok {
		return c.intSliceValues[key]
	}

	return defaultValue
}

func (c *config) String(key string, defaultValue string) string {
	if _, ok := c.stringValues[key]; ok {
		return c.stringValues[key]
	}

	return defaultValue
}

func (c *config) Bool(key string, defaultValue bool) bool {
	if _, ok := c.boolValues[key]; ok {
		return c.boolValues[key]
	}

	return defaultValue
}

func (c *config) Int(key string, defaultValue int) int {
	if _, ok := c.intValues[key]; ok {
		return c.intValues[key]
	}

	return defaultValue
}

func (c *config) WithIntSlice(key string, value []int) contract.Config {
	i := c.clone()
	i.intSliceValues[key] = value

	return i
}

func (c *config) WithString(key, value string) contract.Config {
	i := c.clone()
	i.stringValues[key] = value

	return i
}

func (c *config) WithInt(key string, value int) contract.Config {
	i := c.clone()
	i.intValues[key] = value

	return i
}

func (c *config) WithBool(key string, value bool) contract.Config {
	i := c.clone()
	i.boolValues[key] = value

	return i
}
