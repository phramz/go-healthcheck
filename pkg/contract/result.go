package contract

type Result interface {
	Passed() bool
	Failed() bool
	Reason() string
	Error() error
}
