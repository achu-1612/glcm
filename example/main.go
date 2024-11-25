package main

import "github.com/achu-1612/glcm"

func main() {
	base := glcm.NewRunner()
	base.BootUp(nil)
	base.Wait()
}
