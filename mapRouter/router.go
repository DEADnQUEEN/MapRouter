package mapRouter

import (
	"MapCoder/utils"
	"context"
	"github.com/jackc/pgx/v5"
)

type Router struct {
	loaded  map[int]*Node
	visited map[int]*Node

	conn *pgx.Conn
}

func (r *Router) Clear() {
	for _, node := range r.loaded {
		node.Clear()
	}
	for _, node := range r.visited {
		node.Clear()
	}

	r.loaded = map[int]*Node{}
	r.visited = map[int]*Node{}
}

func (r *Router) GetLoadedNodes() map[int]*Node {
	return r.loaded
}

func (r *Router) findNodeRoads(rootNode *Node) []*Road {
	if rootNode == nil {
		panic("node can't be nil")
	}

	if rootNode.connected != nil {
		rootNode.OverrideLength()
		return rootNode.connected
	}

	r.visited[rootNode.ID] = rootNode

	var args = pgx.NamedArgs{
		"id": rootNode.ID,
	}
	rows, err := r.conn.Query(
		context.Background(),
		"select r.id, c2.* from road r join crossroad c2 on c2.id = r.to_crossroad where r.from_crossroad = @id;",
		args,
	)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var roadId int
		var toNode int
		var toLatitude float64
		var toLongitude float64

		if err = rows.Scan(&roadId, &toNode, &toLongitude, &toLatitude); err != nil {
			panic(err)
		}

		var nodeTo, okTo = r.loaded[toNode]
		if !okTo {
			nodeTo = &Node{ID: toNode, Longitude: toLongitude, Latitude: toLatitude, compareTo: rootNode.compareTo}
			r.loaded[toNode] = nodeTo
		}

		var road = &Road{id: roadId, fromNode: rootNode, toNode: nodeTo}

		rootNode.AddRoad(road)
	}

	rootNode.OverrideLength()
	return rootNode.connected
}

func (r *Router) FindRoute(from *Node, to *Node) float64 {
	if from.ID == to.ID {
		return 0
	}

	r.loaded[from.ID] = from
	r.loaded[to.ID] = to

	from.shortestParent = nil
	from.shortedLength = 0
	from.compareTo = to
	from.isStart = true

	to.compareTo = to

	var nodesQueue = utils.OrderedLinkedList[*Node]{}

	var roads = r.findNodeRoads(from)
	for _, road := range roads {
		road.toNode.shortestParent = from
		nodesQueue.Add(&road.toNode)
	}

	for nodesQueue.GetCount() > 0 {
		var node, err = nodesQueue.RemoveAt(0)
		if err != nil {
			panic(err)
		}

		ok, er := (*node).ExactItem(to)
		if er != nil {
			panic(er)
		}
		if ok {
			return (*node).shortedLength
		}

		for _, road := range r.findNodeRoads(*node) {
			_, visited := r.visited[(*road.toNode).ID]

			if !visited {
				// reorder queue`
				_ = nodesQueue.Remove(&road.toNode)
				nodesQueue.Add(&road.toNode)
			}
		}
	}

	panic("unreachable")
}

func (r *Router) LoadNode(id int) *Node {
	var args = pgx.NamedArgs{"id": id}
	var row = r.conn.QueryRow(
		context.Background(),
		"SELECT longitude, latitude FROM crossroad WHERE id=@id",
		args,
	)

	var node = Node{ID: id}

	err := row.Scan(&node.Longitude, &node.Latitude)
	if err != nil {
		panic(err)
	}

	return &node
}

func CreateRouter(conn *pgx.Conn) *Router {
	return &Router{
		loaded:  map[int]*Node{},
		visited: map[int]*Node{},
		conn:    conn,
	}
}
