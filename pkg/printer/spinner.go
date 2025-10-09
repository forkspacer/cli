package printer

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/forkspacer/cli/pkg/styles"
)

// Spinner wraps the spinner library with our styles
type Spinner struct {
	s *spinner.Spinner
}

// NewSpinner creates a new spinner with default settings
func NewSpinner(message string) *Spinner {
	s := spinner.New(
		spinner.CharSets[14], // Dots
		100*time.Millisecond,
		spinner.WithColor("fgHiCyan"),
	)
	s.Suffix = " " + message + "..."
	return &Spinner{s: s}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	s.s.Start()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.s.Stop()
}

// Success stops the spinner and shows success message
func (s *Spinner) Success(message string) {
	s.s.Stop()
	fmt.Println(styles.Success(message))
}

// Error stops the spinner and shows error message
func (s *Spinner) Error(message string) {
	s.s.Stop()
	fmt.Println(styles.Error(message))
}

// UpdateMessage changes the spinner message
func (s *Spinner) UpdateMessage(message string) {
	s.s.Suffix = " " + message + "..."
}
