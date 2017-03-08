package detection

import "fmt"

func StartDetection(windowSize int, splash float32, input <-chan int, output chan<- int) {
	window := []int{}
	var average float32

	clear := func() {
		average = 0
		window = []int{}
	}

	for count := range input {
		if count != 0 {
			fmt.Printf("New item: %d\n", count)

			// Slide full window
			if len(window) == windowSize {
				first := window[0]
				window = window[1:]
				average -= float32(first) / float32(windowSize)
			}

			// Recalculate average value
			average += float32(count) / float32(windowSize)
			window = append(window, count)

			// Try to detect splash value in full window
			if len(window) == windowSize && average*splash < float32(count) {
				clear()
				output <- count
			}
		}
	}
}
