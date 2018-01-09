package main
import(
	"fmt"
	"time"
)

func main(){
	hoy := time.Now()
		fmt.Println(hoy.Truncate(time.Second).Format(time.RFC3339))
		fmt.Println(hoy.Format(time.RFC3339))
		fmt.Println(hoy)
}