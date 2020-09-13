package pid

import (
	"time"
)

type PID struct {
	p float64
	i float64
	d float64
	f float64

	maxIOutput     float64
	maxError       float64
	errorSum       float64
	maxOutput      float64
	minOutput      float64
	setpoint       float64
	lastActual     float64
	outputRampRate float64
	lastOutput     float64
	outputFilter   float64
	setpointRange  float64

	lastCycle time.Time

	firstRun bool
	reversed bool

	info PIDInfo
}

func NewPID(p, i, d, f float64) *PID {
	out := &PID{
		p:         p,
		i:         i,
		d:         d,
		f:         f,
		firstRun:  true,
		reversed:  false,
		lastCycle: time.Now(),
		info:      PIDInfo{},
	}
	out.checkSigns()
	return out
}

func (pid *PID) SetP(P float64) {
	pid.p = P
	pid.checkSigns()
}

func (pid *PID) SetI(I float64) {
	if pid.i != 0 {
		pid.errorSum = pid.errorSum * pid.i / I
	}
	if pid.maxIOutput != 0 {
		pid.maxError = pid.maxIOutput / I
	}
	pid.i = I
	pid.checkSigns()
}

func (pid *PID) SetD(D float64) {
	pid.d = D
	pid.checkSigns()
}

func (pid *PID) SetPID(p, i, d, f float64) {
	pid.p = p
	pid.i = i
	pid.d = d
	pid.f = f
	pid.checkSigns()
}

func (pid *PID) SetMaxIOutput(max float64) {
	pid.maxIOutput = max
	if pid.i != 0 {
		pid.maxError = pid.maxIOutput / pid.i
	}
}

func (pid *PID) SetOutputLimits(min, max float64) {
	if min > max {
		return
	}
	pid.minOutput = min
	pid.maxOutput = max
	if pid.maxOutput == 0 || pid.maxIOutput > (max-min) {
		pid.SetMaxIOutput(max - min)
	}
}

func (pid *PID) SetDirection(reversed bool) {
	pid.reversed = reversed
}

func (pid *PID) SetSetpoint(sp float64) {
	pid.setpoint = sp
}

func (pid *PID) GetOutput(actual float64) float64 {
	var output float64
	var Poutput float64
	var Ioutput float64
	var Doutpout float64
	var Foutput float64

	setpoint := pid.setpoint

	if pid.setpointRange != 0 {
		setpoint = clamp(setpoint, actual-pid.setpointRange, actual+pid.setpointRange)
	}

	//Calculate error
	pidErr := pid.setpoint - actual

	//F output
	Foutput = pid.f * pid.setpoint

	//P output
	Poutput = pid.p * pidErr

	if pid.firstRun {
		pid.lastActual = actual
		pid.lastOutput = Poutput + Foutput
		pid.firstRun = false
	}

	dt := time.Since(pid.lastCycle).Seconds()

	//Calcuate D
	if dt > 0 {
		dWorking := (actual - pid.lastActual) / dt
		Doutpout = -pid.d * dWorking
	}

	pid.lastActual = actual

	//I term
	Ioutput = pid.i * dt * pid.errorSum
	if pid.maxIOutput != 0 {
		Ioutput = clamp(Ioutput, -pid.maxIOutput, pid.maxIOutput)
	}
	//add up the items
	output = Foutput + Poutput + Ioutput + Doutpout

	if (pid.minOutput != pid.maxOutput) && !bounded(output, pid.minOutput, pid.maxOutput) {
		pid.errorSum = pidErr
	} else if (pid.outputRampRate != 0) && !bounded(output, pid.lastOutput-pid.outputRampRate, pid.lastOutput+pid.outputRampRate) {
		pid.errorSum = pidErr
	} else if pid.maxIOutput != 0 {
		pid.errorSum = clamp(pid.errorSum+pidErr, -pid.maxError, pid.maxError)
	} else {
		pid.errorSum += pidErr
	}

	if pid.outputRampRate != 0 {
		output = clamp(output, pid.lastOutput-pid.outputRampRate, pid.lastOutput+pid.outputRampRate)
	}

	if pid.minOutput != pid.maxOutput {
		output = clamp(output, pid.minOutput, pid.maxOutput)
	}

	if pid.outputFilter != 0 {
		output = pid.lastOutput*pid.outputFilter + output*(1-pid.outputFilter)
	}

	pid.lastCycle = time.Now()
	pid.lastOutput = output
	pid.updateDebugInfo(Poutput, Ioutput, Doutpout, pid.setpoint, actual, output, Foutput, pid.errorSum)
	return output
}

func (pid *PID) Reset() {
	pid.firstRun = true
	pid.errorSum = 0
}

func (pid *PID) SetOutputRampRate(rate float64) {
	pid.outputRampRate = rate
}

func (pid *PID) setSetpointRange(srange float64) {
	pid.setpointRange = srange
}

func (pid *PID) setOutputFilter(str float64) {
	if (str == 0) || (bounded(str, 0, 1)) {
		pid.outputFilter = str
	}
}
