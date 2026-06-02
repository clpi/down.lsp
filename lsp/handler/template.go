package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TemplateVariable represents a variable that can be substituted in templates.
type TemplateVariable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     string `json:"default,omitempty"`
}

// TemplateDefinition represents a reusable document template.
type TemplateDefinition struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Category    string             `json:"category"`
	Body        string             `json:"body"`
	Variables   []TemplateVariable `json:"variables,omitempty"`
	Icon        string             `json:"icon"`
}

// BuiltinTemplates contains the built-in templates.
var BuiltinTemplates = []TemplateDefinition{
	{
		Name:        "note",
		Description: "A basic note",
		Category:    "Basic",
		Icon:        "📝",
		Body: `---
title: {{title}}
date: {{date}}
tags: [{{tags}}]
type: note
---

# {{title}}

{{content}}
`,
		Variables: []TemplateVariable{
			{Name: "title", Description: "Note title", Default: "Untitled"},
			{Name: "date", Description: "Creation date"},
			{Name: "tags", Description: "Comma-separated tags", Default: ""},
			{Name: "content", Description: "Initial content", Default: ""},
		},
	},
	{
		Name:        "daily",
		Description: "Daily journal entry",
		Category:    "Journal",
		Icon:        "📅",
		Body: `---
title: {{date}}
date: {{date}}
type: daily
---

# {{date}} — {{day}}

## Focus

- [ ] {{focus1}}
- [ ] {{focus2}}
- [ ] {{focus3}}

## Notes

{{notes}}

## Reflection

- What went well:
- What could improve:
- Tomorrow's priority:
`,
	},
	{
		Name:        "meeting",
		Description: "Meeting notes",
		Category:    "Work",
		Icon:        "🤝",
		Body: `---
title: "{{title}}"
date: {{date}}
type: meeting
attendees: [{{attendees}}]
---

# {{title}}

**Date**: {{date}}
**Attendees**: {{attendees}}

## Agenda

1. {{agenda1}}
2. {{agenda2}}

## Notes

{{notes}}

## Action Items

- [ ] {{action1}} — @{{assignee1}}
- [ ] {{action2}} — @{{assignee2}}

## Decisions

-
`,
	},
	{
		Name:        "project",
		Description: "Project planning document",
		Category:    "Work",
		Icon:        "📋",
		Body: `---
title: "{{title}}"
date: {{date}}
type: project
status: active
---

# {{title}}

## Overview

{{description}}

## Goals

- [ ] {{goal1}}
- [ ] {{goal2}}
- [ ] {{goal3}}

## Tasks

### Todo
- [ ] {{task1}}
- [ ] {{task2}}

### In Progress

### Done

## Timeline

| Milestone | Date | Status |
| --- | --- | --- |
| {{milestone1}} | {{date1}} | pending |

## Notes

`,
	},
	{
		Name:        "weekly",
		Description: "Weekly review",
		Category:    "Journal",
		Icon:        "📆",
		Body: `---
title: "Week of {{date}}"
date: {{date}}
type: weekly
---

# Week of {{date}}

## Accomplishments

-

## Challenges

-

## Next Week

- [ ]
- [ ]
- [ ]

## Metrics

| Metric | This Week | Last Week |
| --- | --- | --- |
| | | |
`,
	},
	{
		Name:        "reference",
		Description: "Reference/knowledge article",
		Category:    "Knowledge",
		Icon:        "📚",
		Body: `---
title: "{{title}}"
date: {{date}}
type: reference
tags: [{{tags}}]
---

# {{title}}

## Summary

{{summary}}

## Details

{{content}}

## Related

- [[]]

## Sources

-
`,
	},
}

// ExpandTemplate substitutes variables in a template body.
func ExpandTemplate(tmpl *TemplateDefinition, vars map[string]string) string {
	body := tmpl.Body

	// Built-in variables
	now := time.Now()
	builtins := map[string]string{
		"date":     now.Format("2006-01-02"),
		"time":     now.Format("15:04"),
		"datetime": now.Format("2006-01-02 15:04"),
		"day":      now.Format("Monday"),
		"month":    now.Format("January"),
		"year":     now.Format("2006"),
		"week":     fmt.Sprintf("%d", isoWeek(now)),
	}

	// Apply builtins first (as defaults)
	for k, v := range builtins {
		if _, exists := vars[k]; !exists {
			vars[k] = v
		}
	}

	// Apply template defaults for unset variables
	for _, tv := range tmpl.Variables {
		if _, exists := vars[tv.Name]; !exists {
			vars[tv.Name] = tv.Default
		}
	}

	// Substitute all {{variable}} patterns
	for k, v := range vars {
		body = strings.ReplaceAll(body, "{{"+k+"}}", v)
	}

	// Remove any remaining unsubstituted variables
	for {
		start := strings.Index(body, "{{")
		if start < 0 {
			break
		}
		end := strings.Index(body[start:], "}}")
		if end < 0 {
			break
		}
		body = body[:start] + body[start+end+2:]
	}

	return body
}

// ListTemplates returns available templates from builtins + workspace templates dir.
func (s *State) ListTemplates() []TemplateDefinition {
	templates := make([]TemplateDefinition, len(BuiltinTemplates))
	copy(templates, BuiltinTemplates)

	// Scan workspace template directories
	for _, doc := range s.Documents {
		_ = doc // Future: scan for template files
	}

	return templates
}

// CreateFromTemplate creates a new document from a template and writes it.
func (s *State) CreateFromTemplate(templateName string, vars map[string]string, outputDir string) (string, error) {
	var tmpl *TemplateDefinition
	for i := range BuiltinTemplates {
		if BuiltinTemplates[i].Name == templateName {
			tmpl = &BuiltinTemplates[i]
			break
		}
	}
	if tmpl == nil {
		return "", fmt.Errorf("template %q not found", templateName)
	}

	content := ExpandTemplate(tmpl, vars)

	// Determine filename
	title := vars["title"]
	if title == "" {
		title = time.Now().Format("2006-01-02") + "-" + templateName
	}
	filename := slugify(title) + ".md"
	fullPath := filepath.Join(outputDir, filename)

	// Ensure directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return "file://" + fullPath, nil
}

func slugify(s string) string {
	s = strings.ToLower(s)
	var result []byte
	prevDash := false
	for _, b := range []byte(s) {
		if (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') {
			result = append(result, b)
			prevDash = false
		} else if b == ' ' || b == '-' || b == '_' || b == '/' {
			if !prevDash && len(result) > 0 {
				result = append(result, '-')
				prevDash = true
			}
		}
	}
	for len(result) > 0 && result[len(result)-1] == '-' {
		result = result[:len(result)-1]
	}
	return string(result)
}

func isoWeek(t time.Time) int {
	_, week := t.ISOWeek()
	return week
}
