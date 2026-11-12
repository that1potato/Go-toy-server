package main

import "net/http"

type experimentHandler struct{}

func (e experimentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

type Celsius float64
type Fahrenheit float64

// func add() float64 {
//     var c Celsius = 12
//     var f Fahrenheit = 13
//     return c + f
// }

func apply(i int, f func(int) int) int {
	return f(i)
}

func haha() int {
	return apply(1, func(a int) int { return a * 2 })
}

func hahaha() int {
	return apply(1, func(a int) int { return a ^ 2 })
}
