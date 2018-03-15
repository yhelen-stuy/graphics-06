package main

import ()

func main() {
	image := MakeImage(500, 500)
	t := MakeMatrix(4, 4)
	t.Ident()
	e := MakeMatrix(4, 0)
	ParseFile("script", t, e, image)
}
