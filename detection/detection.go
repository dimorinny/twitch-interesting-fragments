package detection

import "fmt"

func StartDetection(
	windowSize int,
	spikeRate float32,
	smoothRate float32,
	input <-chan int,
) <-chan float32 {
	output := make(chan float32)

	window := []int{}
	var average float32

	clear := func() {
		average = 0
		window = []int{}
	}

	checkSplash := func(spikeRate float32, current int) bool {
		return average*spikeRate < float32(current)
	}

	getFragmentRate := func(current int) float32 {
		return float32(current)/average - spikeRate + 1
	}

	go func() {
		for count := range input {
			if count != 0 {
				fmt.Printf("New item: %d ", count)

				// Detect full window
				if len(window) == windowSize {
					if checkSplash(spikeRate, count) {
						// Clear window when spike detected
						clear()
						output <- getFragmentRate(count)
						continue
					} else if smoothRate != 0 && checkSplash(spikeRate/smoothRate, count) {
						// Compensate increasing average value
						count = int(average) + (count-int(average))/int(smoothRate)
					}

					first := window[0]
					window = window[1:]
					average -= float32(first) / float32(windowSize)
				}

				fmt.Printf("Updated item: %d ", count)

				// Recalculate average value
				average += float32(count) / float32(windowSize)
				window = append(window, count)

				fmt.Printf("Current average: %f\n", average)
			}
		}
		close(output)
	}()

	return output
}
