package mapRouter

import (
	"math"
	"strconv"
)

const EarthRadius = 6371000

type Road struct {
	id       int
	fromNode *Node
	toNode   *Node
}

func degreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func CalculateLength(fromLatitude float64, fromLongitude float64, toLatitude float64, toLongitude float64) float64 {
	var radLat1 = degreesToRadians(fromLatitude)
	var radLat2 = degreesToRadians(toLatitude)

	return 2 * EarthRadius * math.Asin(
		math.Sqrt(
			math.Cos(
				radLat1,
			)*math.Cos(
				radLat2,
			)*math.Pow(
				math.Sin((degreesToRadians(toLongitude)-degreesToRadians(fromLongitude))/2), 2,
			)+math.Pow(
				math.Sin((radLat2-radLat1)/2),
				2,
			),
		),
	)
}

func (r *Road) GetLength() float64 {
	return CalculateLength(r.fromNode.Latitude, r.fromNode.Longitude, r.toNode.Latitude, r.toNode.Longitude)
}

func (r *Road) String() string {
	return "RoadID: " + strconv.Itoa(r.id) + "; From: " + r.fromNode.String() + "; To: " + r.toNode.String() + ";"
}
