package main

import (
	"MapCoder/mapRouter"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func main() {
	conn, err := pgx.Connect(context.Background(), "postgresql://testuser:test@localhost:5432/postgres")
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	var router = mapRouter.CreateRouter(conn)
	var start = router.LoadNode(6814066735)
	var end = router.LoadNode(8037233757) //8037233756
	fmt.Println(start.Latitude, start.Longitude, end.Latitude, end.Longitude, mapRouter.CalculateLength(start.Latitude, start.Longitude, end.Latitude, end.Longitude))
	var length = router.FindRoute(start, end)

	fmt.Println("Shortest length -> ", length)
	for i, el := range end.GetPath() {
		fmt.Println(i, " | ", el, " | ", el.Length())
	}
	fmt.Println(len(router.GetLoadedNodes()))
}
