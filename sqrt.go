// Based on https://golang.org/pkg/math/big/#example__sqrt2

package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
)

var prec int64

func init() {
	flag.Int64Var(&prec, "precission", 10000, "Computations will be done with these bits in the mantissa")
}

func Sqrt(n *big.Float) *big.Float {
	steps := int(math.Log2(float64(prec)))

	half := new(big.Float).SetPrec(uint(prec)).SetFloat64(0.5)

	x := new(big.Float).SetPrec(uint(prec)).SetInt64(1)

	t := new(big.Float)

	// Iterate.
	for i := 0; i <= steps; i++ {
		t.Quo(n, x)    // t = 2.0 / x_n
		t.Add(x, t)    // t = x_n + (2.0 / x_n)
		x.Mul(half, t) // x_{n+1} = 0.5 * t
	}

	return x
}

func sqrtHandler(w http.ResponseWriter, r *http.Request) {
	n := r.URL.Path[1:]

	limit, ok := big.NewFloat(0.0).SetString(n)
	if !ok {
		http.Error(w, "NaN", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%s\n", Sqrt(limit).String())
}

func main() {
	flag.Parse()
	http.HandleFunc("/", sqrtHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
