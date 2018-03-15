package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseFile(filename string, t *Matrix, e *Matrix, image *Image) error {
	f, err := os.Open(filename)
	if err != nil {
		return errors.New("Couldn't open file")
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		switch c := strings.TrimSpace(scanner.Text()); c {
		case "line":
			args := getArgs(scanner)
			if err := checkArgCount(args, 6); err != nil {
				fmt.Println(err)
				continue
			}
			fargs := numerize(args)
			e.AddEdge(fargs[0], fargs[1], fargs[2], fargs[3], fargs[4], fargs[5])

		case "ident":
			t.Ident()

		case "scale":
			args := getArgs(scanner)
			if err := checkArgCount(args, 3); err != nil {
				fmt.Println(err)
				continue
			}
			fargs := numerize(args)
			scale := MakeScale(fargs[0], fargs[1], fargs[2])
			t, _ = t.Mult(scale)

		case "move":
			args := getArgs(scanner)
			if err := checkArgCount(args, 3); err != nil {
				fmt.Println(err)
				continue
			}
			fargs := numerize(args)
			translate := MakeTranslate(fargs[0], fargs[1], fargs[2])
			t, _ = t.Mult(translate)

		case "rotate":
			args := getArgs(scanner)
			if err := checkArgCount(args, 2); err != nil {
				fmt.Println(err)
				continue
			}
			// TODO: Error handling
			deg, _ := strconv.ParseFloat(args[1], 64)
			fmt.Printf("rotating %s by %.2f\n", args[0], deg)
			switch args[0] {
			case "x":
				rot := MakeRotX(deg)
				fmt.Println(rot)
				t, _ = t.Mult(rot)
			case "y":
				rot := MakeRotY(deg)
				fmt.Println(rot)
				t, _ = t.Mult(rot)
			case "z":
				rot := MakeRotZ(deg)
				fmt.Println(rot)
				t, _ = t.Mult(rot)
			default:
				// TODO: Error handling
				fmt.Println("Rotate fail")
				continue
			}

		case "circle":
			args := getArgs(scanner)
			if err := checkArgCount(args, 4); err != nil {
				fmt.Println(err)
				continue
			}
			fargs := numerize(args)
			e.AddCircle(fargs[0], fargs[1], fargs[2], fargs[3])

		case "hermite":
			args := getArgs(scanner)
			if err := checkArgCount(args, 8); err != nil {
				fmt.Println(err)
				continue
			}
			fargs := numerize(args)
			err := e.AddHermite(fargs[0], fargs[1], fargs[2], fargs[3], fargs[4], fargs[5], fargs[6], fargs[7], 0.01)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case "bezier":
			args := getArgs(scanner)
			if err := checkArgCount(args, 8); err != nil {
				fmt.Println(err)
				continue
			}
			fargs := numerize(args)
			err := e.AddBezier(fargs[0], fargs[1], fargs[2], fargs[3], fargs[4], fargs[5], fargs[6], fargs[7], 0.01)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case "apply":
			fmt.Println(t)
			// TODO: Error handling
			e, _ = e.Mult(t)

		case "display":
			image.Clear()
			image.DrawLines(e, Color{r: 0, b: 0, g: 0})
			image.Display()

		case "save":
			args := getArgs(scanner)
			if err := checkArgCount(args, 1); err != nil {
				fmt.Println(err)
				continue
			}
			image.Clear()
			image.DrawLines(e, Color{r: 0, b: 0, g: 0})
			image.SavePPM(args[0])

		case "quit":
			break

		default:
			fmt.Printf("Error: Couldn't recognize command %s\n", c)
			continue
		}
	}
	return nil
}

func getArgs(s *bufio.Scanner) []string {
	s.Scan()
	arg := strings.TrimSpace(s.Text())
	return strings.Fields(arg)
}

// Returns error if incorrect number of args, nil otherwise
// TODO: Add funcname for better testing? idk
func checkArgCount(args []string, expected int) error {
	if len(args) != expected {
		return fmt.Errorf("Error: Incorrect # of args. Got: %d, expected: %d\n", len(args), expected)
	}
	return nil
}

// Error handling?
func numerize(args []string) []float64 {
	fargs := make([]float64, len(args))
	for i, val := range args {
		fargs[i], _ = strconv.ParseFloat(val, 64)
	}
	return fargs
}
