package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-steward/internal/server";"github.com/stockyard-dev/stockyard-steward/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9700"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./steward-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("steward: %v",err)};defer db.Close();srv:=server.New(db,server.DefaultLimits())
fmt.Printf("\n  Steward — Self-hosted expense tracker\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("steward: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
