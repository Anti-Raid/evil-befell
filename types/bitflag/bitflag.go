package bitflag

import (
	"fmt"
	"math/big"
)

// BitFlag is a basic bitflag wrapper to manage bitflags decently
type BitFlag struct {
	flagDescriptors map[string]any
	flags           *big.Int
}

// NewBitFlag creates a new BitFlag instance
func NewBitFlag(flagDescriptors map[string]any, initialFlags string) *BitFlag {
	flags, _ := new(big.Int).SetString(initialFlags, 10)
	return &BitFlag{
		flagDescriptors: flagDescriptors,
		flags:           flags,
	}
}

// SetFlag sets the flag to the given value
func (bf *BitFlag) SetFlag(flag interface{}, value bool) error {
	flagValue := new(big.Int)
	for flagKey, descriptor := range bf.flagDescriptors {
		if flagKey == flag || descriptor == flag {
			switch v := descriptor.(type) {
			case string:
				flagValue.SetString(v, 10)
			case int:
				flagValue.SetInt64(int64(v))
			}
			break
		}
	}

	if flagValue.Cmp(big.NewInt(0)) == 0 {
		return fmt.Errorf("unknown flag: %v", flag)
	}

	if value {
		bf.flags.Or(bf.flags, flagValue)
	} else {
		bf.flags.AndNot(bf.flags, flagValue)
	}
	return nil
}

// IsFlagSet returns whether the flag is set
func (bf *BitFlag) IsFlagSet(flag interface{}) bool {
	flagValue := new(big.Int)
	for flagKey, descriptor := range bf.flagDescriptors {
		if flagKey == flag || descriptor == flag {
			switch v := descriptor.(type) {
			case string:
				flagValue.SetString(v, 10)
			case int:
				flagValue.SetInt64(int64(v))
			}
			break
		}
	}

	if flagValue.Cmp(big.NewInt(0)) == 0 {
		switch v := flag.(type) {
		case string:
			flagValue.SetString(v, 10)
		case int:
			flagValue.SetInt64(int64(v))
		}
	}

	return new(big.Int).And(bf.flags, flagValue).Cmp(flagValue) == 0
}

// GetFlags returns the flags
func (bf *BitFlag) GetFlags() *big.Int {
	return new(big.Int).Set(bf.flags)
}

// GetSetFlags returns a list of flags that are set
func (bf *BitFlag) GetSetFlags() map[string]string {
	setFlags := make(map[string]string)
	for flagKey, descriptor := range bf.flagDescriptors {
		flagVal := new(big.Int)
		switch v := descriptor.(type) {
		case string:
			flagVal.SetString(v, 10)
		case int:
			flagVal.SetInt64(int64(v))
		}
		if new(big.Int).And(bf.flags, flagVal).Cmp(flagVal) == 0 {
			setFlags[flagKey] = flagVal.String()
		}
	}
	return setFlags
}

// GetUnsetFlags returns a list of flags that are not set
func (bf *BitFlag) GetUnsetFlags() map[string]string {
	unsetFlags := make(map[string]string)
	for flagKey, descriptor := range bf.flagDescriptors {
		flagVal := new(big.Int)
		switch v := descriptor.(type) {
		case string:
			flagVal.SetString(v, 10)
		case int:
			flagVal.SetInt64(int64(v))
		}
		if new(big.Int).And(bf.flags, flagVal).Cmp(flagVal) != 0 {
			unsetFlags[flagKey] = flagVal.String()
		}
	}
	return unsetFlags
}

// GetFlagKey returns the flag key
func (bf *BitFlag) GetFlagKey(flag interface{}) (string, error) {
	for flagKey, descriptor := range bf.flagDescriptors {
		if flagKey == flag || descriptor == flag {
			return flagKey, nil
		}
	}
	return "", fmt.Errorf("unknown flag: %v", flag)
}

// GetFlagDescriptors returns the flag descriptors
func (bf *BitFlag) GetFlagDescriptors() map[string]interface{} {
	return bf.flagDescriptors
}
