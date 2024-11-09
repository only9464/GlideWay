package gitdorker

import "fmt"

type GitdorkerApp struct{}

// Gitdorker is a method of GitdorkerService that calls the package-level function
func (s *GitdorkerApp) Gitdorker(a, b int) int {
	return Gitdorker(a, b)
}

func Gitdorker(a, b int) int {
	fmt.Println("Gitdorker")
	return a + b
}
