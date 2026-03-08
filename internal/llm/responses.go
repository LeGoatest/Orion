package llm

// Response contains the structured output from an LLM call
type Response struct {
	Steps []string
	Raw   string
}
