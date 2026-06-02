package handler

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ChatMessage represents a message in the AI chat panel.
type ChatMessage struct {
	Role      string    `json:"role"` // "user", "assistant", "system"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	DocURI    string    `json:"doc_uri,omitempty"`
}

// ChatSession tracks an ongoing AI chat conversation.
type ChatSession struct {
	ID       string        `json:"id"`
	Messages []ChatMessage `json:"messages"`
	DocURI   string        `json:"doc_uri,omitempty"`
	Created  time.Time     `json:"created"`
}

// ChatResponse is returned from chat commands.
type ChatResponse struct {
	Message  string `json:"message"`
	Sources  []string `json:"sources,omitempty"`
	Actions  []string `json:"actions,omitempty"`
}

// cmdChat handles the main AI chat interaction.
func (s *State) cmdChat(args []interface{}) (any, error) {
	if s.AI == nil {
		return &ChatResponse{Message: "AI engine not initialized. Set ANTHROPIC_API_KEY or another provider key."}, nil
	}

	if len(args) < 1 {
		return &ChatResponse{Message: "Usage: down.chat <message> [documentURI]"}, nil
	}

	message, ok := args[0].(string)
	if !ok || message == "" {
		return &ChatResponse{Message: "Message must be a non-empty string"}, nil
	}

	var docURI string
	if len(args) > 1 {
		docURI, _ = args[1].(string)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*1e9)
	defer cancel()

	// Build enriched context
	enrichedMessage := s.enrichChatMessage(message, docURI)

	answer, err := s.AI.Query(ctx, docURI, enrichedMessage)
	if err != nil {
		return &ChatResponse{Message: fmt.Sprintf("Chat error: %v", err)}, nil
	}

	// Find related sources
	sources := s.findChatSources(message, docURI)

	// Suggest follow-up actions
	actions := s.suggestChatActions(message, answer, docURI)

	return &ChatResponse{
		Message: answer,
		Sources: sources,
		Actions: actions,
	}, nil
}

// enrichChatMessage adds workspace context to a chat message.
func (s *State) enrichChatMessage(message string, docURI string) string {
	var context strings.Builder
	context.WriteString(message)

	// Add current document context if available
	if docURI != "" {
		if text, ok := s.Documents[docURI]; ok {
			// Add document summary (first 500 chars)
			preview := text
			if len(preview) > 500 {
				preview = preview[:500] + "..."
			}
			context.WriteString("\n\n[Current document context]:\n")
			context.WriteString(preview)
		}

		// Add document type info
		info := s.DetectDocumentType(docURI)
		if info != nil {
			context.WriteString(fmt.Sprintf("\n[Document type: %s, title: %s]", info.Type, info.Title))
		}
	}

	// Add relevant entities from the knowledge graph
	if s.Graph != nil {
		words := strings.Fields(strings.ToLower(message))
		var relatedEntities []string
		for _, word := range words {
			if len(word) < 3 {
				continue
			}
			results := s.Graph.Search(word)
			for _, ent := range results {
				if strings.EqualFold(ent.Name, word) {
					relatedEntities = append(relatedEntities, fmt.Sprintf("%s (%s)", ent.Name, ent.Kind))
				}
			}
		}
		if len(relatedEntities) > 0 && len(relatedEntities) <= 10 {
			context.WriteString("\n[Related entities: " + strings.Join(relatedEntities, ", ") + "]")
		}
	}

	return context.String()
}

// findChatSources identifies documents related to the chat query.
func (s *State) findChatSources(message string, docURI string) []string {
	if s.Graph == nil {
		return nil
	}

	var sources []string
	seen := make(map[string]bool)

	words := strings.Fields(strings.ToLower(message))
	for _, word := range words {
		if len(word) < 3 {
			continue
		}
		results := s.Graph.Search(word)
		for _, ent := range results {
			for _, src := range ent.Sources {
				if src.URI != docURI && !seen[src.URI] {
					sources = append(sources, src.URI)
					seen[src.URI] = true
				}
			}
			if len(sources) >= 5 {
				break
			}
		}
		if len(sources) >= 5 {
			break
		}
	}

	return sources
}

// suggestChatActions suggests follow-up actions based on the conversation.
func (s *State) suggestChatActions(question, answer, docURI string) []string {
	var actions []string
	questionLower := strings.ToLower(question)

	// Suggest creating a note if asking about something new
	if strings.Contains(questionLower, "what is") || strings.Contains(questionLower, "explain") {
		actions = append(actions, "down.template.create:reference — Create a reference note from this answer")
	}

	// Suggest searching for related content
	if strings.Contains(questionLower, "find") || strings.Contains(questionLower, "where") {
		actions = append(actions, "down.knowledge.search — Search the knowledge graph")
	}

	// Suggest expanding into a document
	if len(answer) > 200 {
		actions = append(actions, "down.ai.expand — Expand this into a full document")
	}

	// Suggest related documents
	if docURI != "" {
		actions = append(actions, "down.backlinks — View backlinks for this document")
		actions = append(actions, "down.knowledge.related — Find related documents")
	}

	return actions
}

// cmdChatHistory returns recent chat history.
func (s *State) cmdChatHistory() (any, error) {
	if s.AI == nil {
		return "No AI engine initialized", nil
	}

	// The AI engine maintains conversation history internally
	return "Chat history is maintained per-session. Use down.ai.clear to reset.", nil
}

// cmdChatContext shows what context the AI currently has.
func (s *State) cmdChatContext(args []interface{}) (any, error) {
	if s.AI == nil {
		return "No AI engine initialized", nil
	}

	var docURI string
	if len(args) > 0 {
		docURI, _ = args[0].(string)
	}

	var sb strings.Builder
	sb.WriteString("## AI Chat Context\n\n")

	// Provider info
	sb.WriteString(fmt.Sprintf("**Provider**: %s\n", s.AI.ProviderName()))

	// Document count
	sb.WriteString(fmt.Sprintf("**Open documents**: %d\n", len(s.Documents)))

	// Knowledge graph stats
	if s.Graph != nil {
		sb.WriteString(fmt.Sprintf("**Entities**: %d\n", len(s.Graph.Entities)))
		sb.WriteString(fmt.Sprintf("**Relations**: %d\n", len(s.Graph.Relations)))
	}

	// Current document info
	if docURI != "" {
		info := s.DetectDocumentType(docURI)
		if info != nil {
			sb.WriteString(fmt.Sprintf("\n**Current document**: %s (%s)\n", info.Title, info.Type))
			if len(info.Tags) > 0 {
				sb.WriteString(fmt.Sprintf("**Tags**: %s\n", strings.Join(info.Tags, ", ")))
			}
		}
	}

	return sb.String(), nil
}
