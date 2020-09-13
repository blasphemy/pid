package pid

func (p *PID) checkSigns() {
	if p.reversed { //all vals should be negative in reverse mode
		if p.p > 0 {
			p.p *= -1
		}
		if p.i > 0 {
			p.i *= -1
		}
		if p.d > 0 {
			p.d *= -1
		}
		if p.f > 0 {
			p.f *= -1
		}
	} else {
		//non-reverse make all vals positive
		if p.p < 0 {
			p.p *= -1
		}
		if p.i < 0 {
			p.i *= -1
		}
		if p.d < 0 {
			p.d *= -1
		}
		if p.f < 0 {
			p.f *= -1
		}
	}
}

func clamp(val, min, max float64) float64 {
	if val > max {
		return max
	}
	if val < min {
		return min
	}
	return val
}

func bounded(val, min, max float64) bool {
	if val < min {
		return false
	}
	if val > max {
		return false
	}
	return true
}
