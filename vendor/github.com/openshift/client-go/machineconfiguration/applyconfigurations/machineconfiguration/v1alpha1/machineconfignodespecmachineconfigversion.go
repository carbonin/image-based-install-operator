// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// MachineConfigNodeSpecMachineConfigVersionApplyConfiguration represents a declarative configuration of the MachineConfigNodeSpecMachineConfigVersion type for use
// with apply.
type MachineConfigNodeSpecMachineConfigVersionApplyConfiguration struct {
	Desired *string `json:"desired,omitempty"`
}

// MachineConfigNodeSpecMachineConfigVersionApplyConfiguration constructs a declarative configuration of the MachineConfigNodeSpecMachineConfigVersion type for use with
// apply.
func MachineConfigNodeSpecMachineConfigVersion() *MachineConfigNodeSpecMachineConfigVersionApplyConfiguration {
	return &MachineConfigNodeSpecMachineConfigVersionApplyConfiguration{}
}

// WithDesired sets the Desired field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Desired field is set to the value of the last call.
func (b *MachineConfigNodeSpecMachineConfigVersionApplyConfiguration) WithDesired(value string) *MachineConfigNodeSpecMachineConfigVersionApplyConfiguration {
	b.Desired = &value
	return b
}