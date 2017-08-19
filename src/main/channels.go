package main

import "fmt"

func sum(array []int, ch chan int) {
    sum := 0
    for _, v := range array {
        sum += v
    }
    ch <- sum
}

func main() {
    array := []int{7,1,2,3,4,56,6,7,9}
    ch := make(chan int)
    go sum(array[:len(array)/2], ch)
    go sum(array[len(array)/2:], ch)
    x, y := <-ch, <-ch
    fmt.Println(x, y, x+y)
}
