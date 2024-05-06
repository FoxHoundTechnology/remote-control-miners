package domain

// TODO: validation
// TODO: create custom response objects for domain entities
type Scanner struct {
	StartIP  string
	EndIP    string
	Active   bool
	Location string
}

type ScannerController interface {
	Activate() error
	Deactivate() error
	SetAlert() error
}
