package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

var (
	// window        *glfw.Window
	mu            sync.Mutex
	glfwInitOnce  sync.Once
	overlayActive bool
	clickerActive bool
)

func init() {
	runtime.LockOSThread()
}

func main() {
	fmt.Println("Initializing GLFW...")
	glfwInitOnce.Do(func() {
		if err := glfw.Init(); err != nil {
			log.Fatalln("Failed to initialize GLFW:", err)
		}
	})
	defer glfw.Terminate()

	fmt.Println("Press F12 to enable the overlay and start the click listener...")

	go func() {
		events := hook.Start()
		defer hook.End()

		for e := range events {
			if e.Kind == hook.KeyUp && e.Rawcode == 123 {
				fmt.Println("F12 pressed")
				overlayActive = !overlayActive
				if overlayActive {
					go overlay()
					clickerActive = true
				} else {
					clickerActive = false
				}
			}
			if clickerActive && e.Kind == hook.MouseDown && e.Button == hook.MouseMap["left"] {
				overlayActive = !overlayActive
				x, y := robotgo.Location()
				saveLabel(x, y)
			}
		}
	}()

	select {}
}

func overlay() {

	if !overlayActive {
		return
	}

	fmt.Println("GLFW initialized")

	glfw.WindowHint(glfw.TransparentFramebuffer, glfw.True)
	glfw.WindowHint(glfw.Decorated, glfw.False)
	glfw.WindowHint(glfw.Floating, glfw.True)

	window, err := glfw.CreateWindow(2559, 1439, "Overlay", nil, nil)
	if err != nil {
		log.Fatalln("Failed to create window:", err)
	}
	window.MakeContextCurrent()

	fmt.Println("Window created")

	if err := gl.Init(); err != nil {
		log.Fatalln("Failed to initialize GL:", err)
	}

	fmt.Println("OpenGL initialized")

	fmt.Println("Press F12 to toggle the overlay")

	for !window.ShouldClose() {
		if overlayActive {
			draw()
		} else {
			window.SetShouldClose(true)
			window.Hide()
		}
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func saveLabel(x, y int) {
	fmt.Printf("Click at (%d, %d)\n", x, y)
	fmt.Print("Enter a label for this click: ")
	var label string
	fmt.Scanln(&label)
	if label == "" {
		fmt.Println("No label entered; not saved.")
		return
	}

	filePath := "clicks.json"
	clicks := make(map[string]map[string]int)

	if _, err := os.Stat(filePath); err == nil {
		data, err := os.ReadFile(filePath)
		if err == nil {
			_ = json.Unmarshal(data, &clicks)
		}
	}

	clicks[label] = map[string]int{"x": x, "y": y}

	output, _ := json.MarshalIndent(clicks, "", "  ")
	if err := os.WriteFile(filePath, output, 0644); err != nil {
		fmt.Println("Failed to write JSON:", err)
	} else {
		fmt.Print("\033[H\033[2J") // Clear the terminal
		fmt.Printf("Saved label: %s (%d, %d)\n", label, x, y)
	}

	clickerActive = false

	fmt.Println("Press F12 to enable the overlay and start the click listener...")
}

func draw() {
	gl.ClearColor(0.0, 0.0, 0.0, 0.5) // Semi-transparent black
	gl.Clear(gl.COLOR_BUFFER_BIT)
}
