package pid

import "fmt"

type PIDInfo struct {
	POutput  float64
	IOutput  float64
	DOutput  float64
	FOutput  float64
	Setpoint float64
	Actual   float64
	Output   float64
	ErrSum   float64
}

func (p *PID) updateDebugInfo(pout, iout, dout, setpoint, actual, output, fout, errSum float64) {
	p.info = PIDInfo{
		POutput:  pout,
		IOutput:  iout,
		DOutput:  dout,
		FOutput:  fout,
		Setpoint: setpoint,
		Actual:   actual,
		Output:   output,
		ErrSum:   errSum,
	}
}

func (p *PID) Debug() *PIDInfo {
	return &p.info
}

func (p *PIDInfo) String() string {
	return fmt.Sprintf("pOut=%.2f iOut=%.2f dOut=%.2f fOut=%.2f out=%.2f sp=%.2f actual=%.2f errSum=%.2f", p.POutput, p.IOutput, p.DOutput, p.FOutput, p.Output, p.Setpoint, p.Actual, p.ErrSum)
}
