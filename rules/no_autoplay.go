package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// NoAutoplay checks that media elements don't autoplay.
type NoAutoplay struct{}

func (r *NoAutoplay) Name() string { return RuleNoAutoplay }

func (r *NoAutoplay) Description() string {
	return "media elements should not autoplay (disorienting for users)"
}

func (r *NoAutoplay) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		// Check video and audio elements
		if !n.IsElement("video") && !n.IsElement("audio") {
			return true
		}

		if n.HasAttr("autoplay") {
			// Muted video autoplay is more acceptable (no audio disruption)
			// but still flag it as a warning
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  n.Data + " element has autoplay attribute",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Warning,
			})
		}

		return true
	})

	return results
}
