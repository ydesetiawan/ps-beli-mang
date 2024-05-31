package service

import (
	"math"
	"ps-beli-mang/internal/merchant/model"
	"ps-beli-mang/pkg/errs"
	"time"
)

const (
	EarthRadius    = 6371.0 // Earth radius in kilometers
	SpeedKmPerHour = 40.0
	limitPerKM     = 3.0 // Speed in kilometers per hour
)

// EstimateDeliveryTimeTSP estimates the delivery time for a route with multiple waypoints using TSP
func EstimateDeliveryTimeTSP(merchantsLocation []model.Location, userLocation model.Location) (time.Duration, error) {
	// Include start and end points in the list of waypoints
	allWaypoints := append([]model.Location{}, merchantsLocation...)
	allWaypoints = append(allWaypoints, userLocation)
	if len(merchantsLocation) == 1 {
		return estimateDeliveryTimeMulti(allWaypoints)
	}

	// Find the shortest path using the nearest neighbor algorithm
	shortestPath := nearestNeighbor(allWaypoints)

	// Calculate estimated time for the shortest path
	return estimateDeliveryTimeMulti(shortestPath)
}

// EstimateDeliveryTimeMulti estimates the delivery time for a route with multiple waypoints
func estimateDeliveryTimeMulti(waypoints []model.Location) (time.Duration, error) {
	totalDistance := 0.0

	// Calculate distance and time between consecutive waypoints
	for i := 0; i < len(waypoints)-1; i++ {
		start := waypoints[i]
		end := waypoints[i+1]

		// Calculate distance between consecutive waypoints
		distance := HaversineDistance(start.Lat, start.Long, end.Lat, end.Long)
		totalDistance += distance
	}

	if totalDistance > limitPerKM {
		return 0, errs.NewErrBadRequest("The coordinates is too far (> 3km² in Cartesian coordinate system)")
	}

	// Calculate estimated time in hours
	estimatedTimeHours := totalDistance / SpeedKmPerHour

	// Convert estimated time to duration
	estimatedTime := time.Duration(estimatedTimeHours * float64(time.Hour))

	return estimatedTime, nil
}

// haversineDistance calculates the distance between two points using the Haversine formula
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert latitude and longitude from degrees to radians
	lat1Rad := degToRad(lat1)
	lon1Rad := degToRad(lon1)
	lat2Rad := degToRad(lat2)
	lon2Rad := degToRad(lon2)

	// Calculate the difference between latitudes and longitudes
	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	// Calculate the Haversine distance
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := EarthRadius * c

	return distance
}

// degToRad converts degrees to radians
func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180.0)
}

// nearestNeighbor finds the shortest path using the nearest neighbor algorithm
func nearestNeighbor(waypoints []model.Location) []model.Location {
	n := len(waypoints)
	visited := make([]bool, n)
	path := make([]model.Location, 0)

	// Start from the first waypoint
	current := waypoints[0]
	visited[0] = true
	path = append(path, current)

	// Visit each waypoint using the nearest neighbor
	for len(path) < n {
		minDistance := math.Inf(1)
		nearestIndex := -1

		for i, location := range waypoints {
			if !visited[i] {
				distance := HaversineDistance(current.Lat, current.Long, location.Lat, location.Long)
				if distance < minDistance {
					minDistance = distance
					nearestIndex = i
				}
			}
		}

		if nearestIndex != -1 {
			current = waypoints[nearestIndex]
			visited[nearestIndex] = true
			path = append(path, current)
		}
	}

	return path
}
