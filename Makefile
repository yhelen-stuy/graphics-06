all:
	go run draw.go image.go main.go matrix.go parser.go

clean:
	rm *.ppm
