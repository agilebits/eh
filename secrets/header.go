package secrets

// Header allows access to the EHCL entry
type Header struct {
	EHCL EHCL
}

// EHCL is a special entry in the .hcl file that defines encryption parameters
type EHCL struct {
	Encrypted bool
	Key       string

	Service ServiceParams
	Protect []string
}

// ServiceParams is a part of the EHCL entry with crypto service type and parameters
type ServiceParams struct {
	Type      string
	Region    string
	MasterKey string
}
