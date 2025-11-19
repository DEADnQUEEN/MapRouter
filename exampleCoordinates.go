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

	var length, end = router.FindRouteFromCoordinates(43.8001654, 39.4672131, 43.7983596, 39.4695239)

	fmt.Println("Shortest length -> ", length)
	for i, el := range end.GetPath() {
		fmt.Println(i, " | ", el, " | ", el.Length())
	}
	fmt.Println(len(router.GetLoadedNodes()))
}
