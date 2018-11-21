package main

import "fmt"

// http://westerly-lzh.github.io/cn/2014/05/Exponential-Smoothing/

var (
	y = []int64{362, 385, 432, 341, 382, 409,
		498, 387, 473, 513, 582, 474, 544, 582, 681, 557, 628, 707, 773, 592, 627, 725, 854, 661}
	period = int64(6)
	m      = int64(4)

	alpha = float64(0.5)
	beta  = float64(0.4)
	gamma = float64(0.6)
)

func main() {
	prediction := forecast(y, alpha, beta, gamma, period, m)
	fmt.Printf("%+v", prediction)
}

func forecast(y []int64, alpha, beta, gamma float64, period, m int64) []float64 {

	isPass := validateArguments(y, alpha, beta, gamma, period, m)
	if !isPass {
		return nil
	}

	seasons := int64(len(y)) / period
	a0 := calculateInitialLevel(y)
	b0 := calculateInitialTrend(y, period)
	initialSeasonalIndices := calculateSeasonalIndices(y, period, seasons)

	forecast := calculateHoltWinters(y, a0, b0, alpha, beta, gamma, initialSeasonalIndices, period, m)
	return forecast
}

func validateArguments(y []int64, alpha, beta, gamma float64, period, m int64) bool {
	if len(y) == 0 {
		fmt.Println("lengths of y must be > 0 .")
		return false
	}
	if m <= 0 {
		fmt.Println("Value of m must be > 0 .")
		return false
	}
	if m > period {
		fmt.Println("Value of m must be <= period.")
		return false
	}

	if (alpha < 0.0) || (alpha > 1.0) {
		fmt.Println("Value of Alpha should satisfy 0.0 <= alpha <= 1.0")
		return false
	}

	if (beta < 0.0) || (beta > 1.0) {
		fmt.Println("Value of Beta should satisfy 0.0 <= beta <= 1.0")
		return false
	}

	if (gamma < 0.0) || (gamma > 1.0) {
		fmt.Println("Value of Gamma should satisfy 0.0 <= gamma <= 1.0")
		return false
	}
	return true
}

func calculateHoltWinters(y []int64, a0, b0, alpha, beta, gamma float64, initialSeasonalIndices []float64, period, m int64) []float64 {

	St := make([]float64, len(y))
	Bt := make([]float64, len(y))  // 数据的趋势序列
	It := make([]float64, len(y))  // 季节修正因子序列
	Ft := make([]float64, int64(len(y))+m)

	// Initialize base values
	St[1] = a0
	Bt[1] = b0

	for i := int64(0); i < period; i++ {
		It[i] = initialSeasonalIndices[i]
	}

	// Start calculations
	for i := int64(2); i < int64(len(y)); i++ {
		// Calculate overall smoothing
		if i-period >= 0 {
			St[i] = alpha*float64(y[i])/It[i-period] + (1.0-alpha)*(St[i-1]+Bt[i-1])
		} else {
			St[i] = alpha*float64(y[i]) + (1.0-alpha)*(St[i-1]+Bt[i-1])
		}
		// Calculate trend smoothing
		Bt[i] = beta*(St[i]-St[i-1]) + (1-beta)*Bt[i-1]
		// Calculate seasonal smoothing
		if (i - period) >= 0 {
			It[i] = gamma*float64(y[i])/St[i] + (1.0-gamma)*It[i-period]
		}
		// Calculate forecast
		if (i + m) >= period {
			Ft[i+m] = (St[i] + (float64(m) * Bt[i])) * It[i-period+m]
		}
	}
	return Ft
}

func calculateInitialLevel(y []int64) float64 {
	return float64(y[0])
}

func calculateInitialTrend(y []int64, period int64) float64 {

	sum := float64(0)

	for i := int64(0); i < period; i++ {
		sum += float64(y[period+i] - y[i])
	}

	return sum / float64(period*period)
}

func calculateSeasonalIndices(y []int64, period, seasons int64) []float64 {

	seasonalAverage := make([]float64, seasons)
	seasonalIndices := make([]float64, period)
	averagedObservations := make([]float64, len(y))

	for i := int64(0); i < seasons; i++ {
		for j := int64(0); j < period; j++ {
			seasonalAverage[i] += float64(y[(i*period)+j])
		}
		seasonalAverage[i] /= float64(period)
	}
	for i := int64(0); i < seasons; i++ {
		for j := int64(0); j < period; j++ {
			averagedObservations[(i*period)+j] = float64(y[(i*period)+j]) / seasonalAverage[i]
		}
	}

	for i := int64(0); i < period; i++ {
		for j := int64(0); j < seasons; j++ {
			seasonalIndices[i] += averagedObservations[(j*period)+i]
		}
		seasonalIndices[i] /= float64(seasons)
	}

	return seasonalIndices
}

