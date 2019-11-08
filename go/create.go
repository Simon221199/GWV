package main

import (
	"fmt"
)

func createEnv(width, height int) {

	env := ""

	for iny := 0; iny < height; iny++ {
		for inx := 0; inx < width; inx++ {

			border := inx == 0 || inx == width-1
			border = border || iny == 0 || iny == height-1

			if border {
				env += "x"
				continue
			}

			// blocked := rand.Intn(100) > 82
			//
			// if blocked {
			// 	env += "x"
			// 	continue
			// }

			env += " "
		}

		env += "\n"
	}

	fmt.Println(env)
}
