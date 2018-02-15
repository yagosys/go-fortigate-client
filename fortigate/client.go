// WARNING: This file was generated by gen/generator.go

package fortigate

// A fortigate API client
type Client interface {

	// List all FirewallAddresss
	ListFirewallAddresss() ([]FirewallAddress, error)

	// Get a FirewallAddress by name
	GetFirewallAddress(name string) (FirewallAddress, error)

	// Create a new FirewallAddress
	CreateFirewallAddress(*FirewallAddress) error

	// Update a FirewallAddress
	UpdateFirewallAddress(*FirewallAddress) error

	// Delete a FirewallAddress by name
	DeleteFirewallAddress(name string) error

	// List all FirewallPolicys
	ListFirewallPolicys() ([]FirewallPolicy, error)

	// Get a FirewallPolicy by name
	GetFirewallPolicy(name string) (FirewallPolicy, error)

	// Create a new FirewallPolicy
	CreateFirewallPolicy(*FirewallPolicy) error

	// Update a FirewallPolicy
	UpdateFirewallPolicy(*FirewallPolicy) error

	// Delete a FirewallPolicy by name
	DeleteFirewallPolicy(name string) error

	// List all VIPs
	ListVIPs() ([]VIP, error)

	// Get a VIP by name
	GetVIP(name string) (VIP, error)

	// Create a new VIP
	CreateVIP(*VIP) error

	// Update a VIP
	UpdateVIP(*VIP) error

	// Delete a VIP by name
	DeleteVIP(name string) error
}
