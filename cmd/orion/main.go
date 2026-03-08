package main
import "orion/internal/runtime"
import "fmt"
func main() {
	k, _ := runtime.NewKernel("data")
	k.Start()
	fmt.Println("Orion is running.")
	select {}
}
