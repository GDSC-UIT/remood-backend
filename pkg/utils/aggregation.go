package utils 

// Weight number formula: weight = 0.0036*point^2 - 0.36*point + 10
func GetWeight(point int) float32 {
	var a float32 = 0.0036
	var b float32 = -0.36
	var c float32 = 10
	
	weight := a*float32(point)*float32(point) + b*float32(point) + c 
	return weight
}