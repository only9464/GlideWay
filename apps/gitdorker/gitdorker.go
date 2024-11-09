package gitdorker

import "fmt"

type App struct{}

// Gitdorker is a method of GitdorkerService that calls the package-level function
func (s *App) Gitdorker(a, b int) int {
	return Gitdorker(a, b)
}

func Gitdorker(a, b int) int {
	fmt.Println("Gitdorker")
	return a + b
}
