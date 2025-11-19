This is a project to routing using A* algorithm to provide fast calculate for a fasted route

# Requirements:
## Database:
> [!Note]
> This tables and columns are required for db
- Crossroad table:
    - Longitude
    - Latitude
- Road table
    - Road id
    - From crossroad
    - To crossroad
# Examples
## Code Example:
```go 
conn, err := pgx.Connect(context.Background(), "postgresql://testuser:test@localhost:5432/postgres")
if err != nil {  
    panic(err)  
}  
defer conn.Close(context.Background())  
  
var router = mapRouter.CreateRouter(conn)  

var start = router.LoadNode(from node id)  
var end = router.LoadNode(to node id)

var length = router.FindRoute(start, end)  
  
fmt.Println("Shortest length -> ", length)  
for i, el := range end.GetPath() {  
    fmt.Println(i, " | ", el, " | ", el.Length())  
}  
fmt.Println(len(router.GetLoadedNodes()))
```
## DB view example:
### Crossroad
```csv
id,longitude,latitude  
6920440696,39.0431893,45.334922  
6920440697,39.045901,45.3352727  
6920440698,39.0553949,45.3378687  
6920440699,39.0660242,45.333971  
6920440700,39.0663308,45.3337555
```
### Road
```csv
id,from_crossroad,to_crossroad  
1241664,12775699097,12775699098  
1241665,12775699098,12775699099  
1241666,12775699099,12775699100  
1241667,12775699100,12775701901  
1241668,12775701901,12775701902
```
