package config

import (
	"fmt"
	"regexp"
)

// Validate validates the entire QuzzConfig
func (c *QuzzConfig) Validate() error {
	if err := c.validateLogLevel(); err != nil {
		return err
	}

	if err := c.validateOutputFormat(); err != nil {
		return err
	}

	if err := c.Performance.Validate(); err != nil {
		return fmt.Errorf("performance config: %w", err)
	}

	if err := c.Props.Validate(); err != nil {
		return fmt.Errorf("props config: %w", err)
	}

	if err := c.validateComponentFilter(); err != nil {
		return err
	}

	if err := c.validateThrottleMs(); err != nil {
		return err
	}

	if err := c.validateMaxValues(); err != nil {
		return err
	}

	return nil
}

// validateLogLevel checks if the log level is valid
func (c *QuzzConfig) validateLogLevel() error {
	validLevels := map[string]bool{
		"error": true,
		"warn":  true,
		"info":  true,
		"debug": true,
	}

	if !validLevels[c.LogLevel] {
		return fmt.Errorf("invalid log level '%s': must be one of error, warn, info, debug", c.LogLevel)
	}

	return nil
}

// validateOutputFormat checks if the output format is valid
func (c *QuzzConfig) validateOutputFormat() error {
	validFormats := map[string]bool{
		"pretty":  true,
		"compact": true,
		"json":    true,
	}

	if !validFormats[c.OutputFormat] {
		return fmt.Errorf("invalid output format '%s': must be one of pretty, compact, json", c.OutputFormat)
	}

	return nil
}

// validateComponentFilter checks if the component filter regex is valid
func (c *QuzzConfig) validateComponentFilter() error {
	if c.ComponentFilter == "" {
		return nil
	}

	_, err := regexp.Compile(c.ComponentFilter)
	if err != nil {
		return fmt.Errorf("invalid component filter regex: %w", err)
	}

	return nil
}

// validateThrottleMs checks if throttle value is reasonable
func (c *QuzzConfig) validateThrottleMs() error {
	if c.ThrottleMs < 0 {
		return fmt.Errorf("throttleMs cannot be negative: got %d", c.ThrottleMs)
	}

	if c.ThrottleMs > 10000 {
		return fmt.Errorf("throttleMs too large: got %d (max 10000ms)", c.ThrottleMs)
	}

	return nil
}

// validateMaxValues checks if max depth and length values are reasonable
func (c *QuzzConfig) validateMaxValues() error {
	if c.MaxPropDepth < 0 {
		return fmt.Errorf("maxPropDepth cannot be negative: got %d", c.MaxPropDepth)
	}

	if c.MaxPropDepth > 100 {
		return fmt.Errorf("maxPropDepth too large: got %d (max 100)", c.MaxPropDepth)
	}

	if c.MaxStringLength < 0 {
		return fmt.Errorf("maxStringLength cannot be negative: got %d", c.MaxStringLength)
	}

	if c.MaxStringLength > 100000 {
		return fmt.Errorf("maxStringLength too large: got %d (max 100000)", c.MaxStringLength)
	}

	return nil
}

// Validate validates the PerformanceConfig
func (p *PerformanceConfig) Validate() error {
	if p.WarnThreshold < 0 {
		return fmt.Errorf("warnThreshold cannot be negative: got %d", p.WarnThreshold)
	}

	if p.WarnThreshold > 60000 {
		return fmt.Errorf("warnThreshold too large: got %d (max 60000ms)", p.WarnThreshold)
	}

	if p.MemoryThreshold < 0 {
		return fmt.Errorf("memoryThreshold cannot be negative: got %d", p.MemoryThreshold)
	}

	return nil
}

// Validate validates the PropsConfig
func (p *PropsConfig) Validate() error {
	if p.AwaitTimeout < 0 {
		return fmt.Errorf("awaitTimeout cannot be negative: got %d", p.AwaitTimeout)
	}

	if p.AwaitTimeout > 30000 {
		return fmt.Errorf("awaitTimeout too large: got %d (max 30000ms)", p.AwaitTimeout)
	}

	if p.MaxArrayItems < 0 {
		return fmt.Errorf("maxArrayItems cannot be negative: got %d", p.MaxArrayItems)
	}

	if p.MaxObjectProps < 0 {
		return fmt.Errorf("maxObjectProps cannot be negative: got %d", p.MaxObjectProps)
	}

	if p.MaxErrorDepth < 0 {
		return fmt.Errorf("maxErrorDepth cannot be negative: got %d", p.MaxErrorDepth)
	}

	validStrategies := map[string]bool{
		"basic":  true,
		"deep":   true,
		"custom": true,
	}

	if p.SerializationStrategy != "" && !validStrategies[p.SerializationStrategy] {
		return fmt.Errorf("invalid serialization strategy '%s': must be one of basic, deep, custom", p.SerializationStrategy)
	}

	return nil
}
