package scriba

type StepFunc func() (error, string)

type Step struct {
	Desc string
	Help string
	Func StepFunc
}
