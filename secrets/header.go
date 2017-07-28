package secrets

// Wrapper allows access to the echl Header entry
type Wrapper struct {
	Header Header `hcl:"eh"`
}

// Header is a special entry in the .hcl file that defines encryption parameters
type Header struct {
	Encrypted bool
	Key       string

	Service ServiceParams
	Protect []string
	Include []string
}

// ServiceParams is a part of the header entry with crypto service type and parameters
type ServiceParams struct {
	Type      string
	Region    string
	MasterKey string
}
