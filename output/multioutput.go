package output

type multiOutput struct {
	members []Outputer
}

// NewMultiOutput constructs multiOutput object
func NewMultiOutput(outers ...Outputer) (Outputer, error) {
	return &multiOutput{outers}, nil
}

func (mo multiOutput) Write(p []byte) (n int, err error) {
	for _, o := range mo.members {
		o.Write(p)
	}
	return len(p), nil
}

func (mo *multiOutput) Close() {
	for _, o := range mo.members {
		o.Close()
	}
}
