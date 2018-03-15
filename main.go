package main

import ()

func main() {
	image := MakeImage(500, 500)
	t := MakeMatrix(4, 4)
	t.Ident()
	e := MakeMatrix(4, 0)
	ParseFile("script1", t, e, image)
	// ParseFile("script", t, e, image)
	// t.Ident()
	// e = MakeMatrix(4, 0)
	// ParseFile("darthvaderscript", t, e, image)
}
