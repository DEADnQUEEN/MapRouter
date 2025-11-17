package mapRouter

import (
	"MapCoder/utils"
	"errors"
	"fmt"
	"math"
)

type Node struct {
	ID        int
	Longitude float64
	Latitude  float64

	shortestParent *Node
	shortedLength  float64

	isStart     bool
	compareTo   *Node
	parentCount int

	roadsLoaded bool
	connected   []*Road
}

func (n *Node) Length() float64 {
	return n.shortedLength
}

func (n *Node) Clear() {
	n.shortestParent = nil
	n.shortedLength = 0
}

func (n *Node) PathFunction(node *Node) float64 {
	return math.Sqrt(
		math.Pow(
			(n.Longitude-node.Longitude)*math.Pow(10, 7),
			2,
		)+math.Pow(
			(n.Latitude-node.Latitude)*math.Pow(10, 7),
			2,
		),
	) + (n.shortedLength * float64(n.parentCount))
}

func (n *Node) OverrideLength() {
	for _, road := range n.connected {
		var newLength = n.shortedLength + road.GetLength()
		if (road.toNode.shortedLength > newLength || road.toNode.shortestParent == nil) && !road.toNode.isStart {
			road.toNode.shortedLength = newLength
			road.toNode.parentCount = n.parentCount + 1
			road.toNode.shortestParent = n
			road.toNode.OverrideLength()
		}
	}
}

func (n *Node) GetPath() []*Node {
	if n.shortestParent == nil {
		return []*Node{}
	}

	return append(n.shortestParent.GetPath(), n)
}

func (n *Node) CompareValues(c utils.Comparable) (int, error) {
	switch c.(type) {
	case *Node:
		var compareNode = (c.(*Node)).PathFunction(n.compareTo)
		var currentNode = n.PathFunction(n.compareTo)

		if compareNode < currentNode {
			return utils.LESS, nil
		}
		if compareNode > currentNode {
			return utils.GREATER, nil
		}

		return utils.EQUAL, nil
	}
	return 0, errors.New("not a node")
}

func (n *Node) ExactItem(c utils.Comparable) (bool, error) {
	switch c.(type) {
	case *Node:
		return (c.(*Node)).ID == n.ID, nil
	}

	return false, errors.New("not a node")
}

func (n *Node) String() string {
	return fmt.Sprintf("ID: %d", n.ID)
}

func (n *Node) GetAbsoluteLengthToNode(node *Node) float64 {
	var radLat1 = degreesToRadians(n.Latitude)
	var radLon1 = degreesToRadians(n.Longitude)
	var radLat2 = degreesToRadians(node.Latitude)
	var radLon2 = degreesToRadians(node.Longitude)

	return 2 * EarthRadius * math.Asin(
		math.Sqrt(
			math.Cos(
				radLat1,
			)*math.Cos(
				radLat2,
			)*math.Pow(
				math.Sin((radLon2-radLon1)/2), 2,
			)+math.Pow(
				math.Sin((radLat2-radLat1)/2),
				2,
			),
		),
	)
}

func (n *Node) GetDegreeDeltaForDistance(distance float64, digitsAfterDot int) float64 {
	var step = 180.0
	var checkNode = &Node{Longitude: n.Longitude, Latitude: n.Latitude}

	var length = n.GetAbsoluteLengthToNode(checkNode)
	var maxLength = distance + math.Pow(0.1, float64(digitsAfterDot))

	for length > maxLength || length < distance {
		for length > maxLength {
			checkNode.Latitude -= step
			checkNode.Longitude -= step
			length = n.GetAbsoluteLengthToNode(checkNode)
		}
		for length < distance {
			checkNode.Latitude += step
			checkNode.Longitude += step
			length = n.GetAbsoluteLengthToNode(checkNode)
		}
		step /= 2
	}
	return checkNode.Latitude - n.Latitude
}

func (n *Node) AddRoad(road *Road) {
	if n.connected == nil {
		n.connected = make([]*Road, 0, 1)
	}
	n.connected = append(n.connected, road)
}
