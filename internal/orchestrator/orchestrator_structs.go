package orchestrator

type outputData interface {
	GetData() string
}

type successOutputData struct {
	Result string `json:"result"`
}

type failureOutputData struct {
	Error string `json:"error"`
}

type inputData struct {
	Expression string `json:"expression"`
}

// successOutputData methods
func (sod successOutputData) GetData() string {
	return sod.Result
}

// failureOutputData methods
func (fod failureOutputData) GetData() string {
	return fod.Error
}
