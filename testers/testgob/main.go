package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math"
)

type Point struct {
	X, Y int
}

func (p Point) Hypotenuse() float64 {
	return math.Hypot(float64(p.X), float64(p.Y))
}

type Pythagoras interface {
	Hypotenuse() float64
}

// This shows how to encode the value of an interface type
// The difference between a key and a regular type is to register the specific type of the implementation interface.
func main() {
	var network bytes.Buffer // Standard input
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(Point{})
	// Create an encoder interface and send values
	enc := gob.NewEncoder(&network)
	for i := 1; i <= 3; i++ {
		interfaceEncode(enc, Point{3 * i, 4 * i})
	}
	// Create a decoder interface and receive values
	dec := gob.NewDecoder(&network)
	for i := 1; i <= 3; i++ {
		result := interfaceDecode(dec)
		fmt.Println(result.Hypotenuse())
	}
}

// interfaceEncode encodes and saves the value to the encoder
func interfaceEncode(enc *gob.Encoder, p Pythagoras) {
	// The encode will fail unless the concrete type has been
	// registered. We registered it in the calling function.
	// Pass pointer to interface so Encode sees (and hence sends) a value of
	// interface type.  If we passed p directly it would see the concrete type instead.
	// See the blog post, "The Laws of Reflection" for background.
	err := enc.Encode(&p)
	if err != nil {
		log.Fatal("encode:", err)
	}
}

// interfaceDecode decodes the value of the interface and returns
func interfaceDecode(dec *gob.Decoder) Pythagoras {
	// The decode will fail unless the concrete type on the wire has been
	// registered. We registered it in the calling function.
	var p Pythagoras
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal("decode:", err)
	}
	return p
}
