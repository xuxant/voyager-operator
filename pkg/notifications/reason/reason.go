package reason

const (
	OperatorSource   Source = "operator"
	KubernetesSource Source = "kubernetes"
	HumanSource      Source = "human"
)

type Reason interface {
	Short() []string
	Verbose() []string
	HasMessages() bool
}

type Undefined struct {
	source  Source
	short   []string
	verbose []string
}

type BaseConfigurationFailed struct {
	Undefined
}

type ReconcileLoopFailed struct {
	Undefined
}

type Source string

func (p Undefined) Short() []string {
	return p.short
}

func (p Undefined) Verbose() []string {
	return p.verbose
}

func (p Undefined) HasMessages() bool {
	return len(p.short) > 0 || len(p.verbose) > 0
}

func NewBaseConfigurationFailed(source Source, short []string, verbose ...string) *BaseConfigurationFailed {
	return &BaseConfigurationFailed{
		Undefined{
			source:  source,
			short:   short,
			verbose: checkIfVerboseEmpty(short, verbose),
		},
	}
}

func NewReconcileLoopFailed(source Source, short []string, verbose ...string) *ReconcileLoopFailed {
	return &ReconcileLoopFailed{
		Undefined{
			source:  source,
			short:   short,
			verbose: checkIfVerboseEmpty(short, verbose),
		},
	}
}

func checkIfVerboseEmpty(short []string, verbose []string) []string {
	if len(verbose) == 0 {
		return short
	}

	return verbose
}
