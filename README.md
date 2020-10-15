# DHT Ring Density

Generate a random set of identifiers on a DHT ring, and then determine the ideal
location of a new identifier on the ring. The "ideal" location is defined as 
in the largest gap between two existing identifiers that are in the longest 
contiguous region of lowest ring density.

To build the code:

  go build cmd/density/main.go
