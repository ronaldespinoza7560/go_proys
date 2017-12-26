package main
// consideramos tres estados
// primer estado generar
import (
	"fmt"
)
func gen(nums []int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}
// The second stage, sq, receives integers from a channel and 
// returns a channel that emits the square of each received integer.
//  After the inbound channel is closed and this stage has sent 
// all the values downstream, it closes the outbound channel:

func sq(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func main() {
	// Set up the pipeline.
	arr:=[]int{1,2,3,4,5,6}
    c := gen(arr)
	out := sq(sq(sq(c)))
	n:=0

	// Consume the output.
	for n < len(arr){
		fmt.Println(<-out)
		n++
	}
    
}