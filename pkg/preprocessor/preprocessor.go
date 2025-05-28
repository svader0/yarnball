package preprocessor

import (
	"strings"
)

// The preprocessor is responsible for cleaning up the input Yarnball code.
// There are a lot of little things that make Yarnball code look more like crochet,
// but are not, in fact, needed for the code to work at all (as of right now).

type Preprocessor struct{}

func New() *Preprocessor {
	return &Preprocessor{}
}

func (p *Preprocessor) Process(input string) (string, error) {
	lines := strings.Split(input, "\n")
	var processedLines []string

	// Remove everything before and including the line that says "STITCH GUIDE:" (case-insensitive).
	var startIdx int
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.EqualFold(trimmed, "STITCH GUIDE:") {
			startIdx = i + 1 // Start processing from the next line
			break
		}
		if strings.EqualFold(trimmed, "INSTRUCTIONS:") {
			startIdx = i + 1 // Start processing from the next line
			break
		}
	}

	for _, line := range lines[startIdx:] {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue // Skip empty lines and comment lines
		}
		// If the line says "INSTRUCTIONS:", skip it --- that's just our label to separate the stitch guide (functions) from the actual instructions.
		if strings.EqualFold(strings.TrimSpace(line), "INSTRUCTIONS:") {
			continue
		}

		// Remove leading and trailing whitespace
		line = strings.TrimSpace(line)

		// Remove any comments from the line
		line = p.removeComment(line)

		// Remove "Row #:" or "Round #:" prefixes if present
		line = p.RemoveRowRoundPrefix(line)

		processedLines = append(processedLines, line)
	}

	return strings.Join(processedLines, "\n"), nil
}

func (p *Preprocessor) removeComment(line string) string {
	// Check if the line contains a comment
	if idx := strings.Index(line, "#"); idx != -1 {
		return strings.TrimSpace(line[:idx]) // Return the line up to the comment
	}
	return line
}

// RemoveRowRoundPrefix removes the "Row N:" or "Round N:" prefix from a line if it exists.
// We don't need those, they just add a little crochet-inspired flair to the code.
func (p *Preprocessor) RemoveRowRoundPrefix(line string) string {
	if strings.HasPrefix(line, "Row ") || strings.HasPrefix(line, "Round ") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) > 1 {
			return strings.TrimSpace(parts[1]) // Return the part after the colon
		}
	}
	return line
}
