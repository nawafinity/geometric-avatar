package main

import (
	"fmt"
	"github.com/nawafinity/geometric-avatar"
	"image/color"
)

func main() {
	email := "geometric-avatar"
	options := map[string]interface{}{
		"size":       128.0,
		"margin":     0.08,
		"background": color.RGBA{238, 238, 238, 255},
	}

	avatar := geometric_avatar.GenerateAvatar(email, options)
	dc := avatar.Render()

	// Save the avatar as a PNG file
	err := dc.SavePNG("output.png")
	if err != nil {
		fmt.Println("Error:", err)
	}

}
