// Copyright 2025 Redpanda Data, Inc.

package tracer

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	yaml "gopkg.in/yaml.v3"

	"github.com/redpanda-data/benthos/v4/internal/docs"
)

func init() {
	// TODO: I'm so confused, these APIs are a nightmare.
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

// Config is the all encompassing configuration struct for all tracer types.
type Config struct {
	Type   string `json:"type" yaml:"type"`
	Plugin any    `json:"plugin,omitempty" yaml:"plugin,omitempty"`
}

// NewConfig returns a configuration struct fully populated with default values.
func NewConfig() Config {
	return Config{
		Type:   "none",
		Plugin: nil,
	}
}

// FromAny returns a tracer config from a parsed config, yaml node or map.
func FromAny(prov docs.Provider, value any) (conf Config, err error) {
	switch t := value.(type) {
	case Config:
		return t, nil
	case *yaml.Node:
		return fromYAML(prov, t)
	case map[string]any:
		return fromMap(prov, t)
	}
	err = fmt.Errorf("unexpected value, expected object, got %T", value)
	return
}

func fromMap(prov docs.Provider, value map[string]any) (conf Config, err error) {
	if conf.Type, _, err = docs.GetInferenceCandidateFromMap(prov, docs.TypeTracer, value); err != nil {
		err = docs.NewLintError(0, docs.LintComponentNotFound, err)
		return
	}

	if p, exists := value[conf.Type]; exists {
		conf.Plugin = p
	} else if p, exists := value["plugin"]; exists {
		conf.Plugin = p
	}
	return
}

func fromYAML(prov docs.Provider, value *yaml.Node) (conf Config, err error) {
	if conf.Type, _, err = docs.GetInferenceCandidateFromYAML(prov, docs.TypeTracer, value); err != nil {
		err = docs.NewLintError(value.Line, docs.LintComponentNotFound, err)
		return
	}

	pluginNode, err := docs.GetPluginConfigYAML(conf.Type, value)
	if err != nil {
		err = docs.NewLintError(value.Line, docs.LintFailedRead, err)
		return
	}

	conf.Plugin = &pluginNode
	return
}
