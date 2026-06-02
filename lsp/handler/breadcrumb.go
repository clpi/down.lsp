package handler

import (
	"path/filepath"
	"strings"
)

// DocumentType represents the detected type of a markdown document.
type DocumentType string

const (
	DocTypeNote      DocumentType = "note"
	DocTypeDaily     DocumentType = "daily"
	DocTypeWeekly    DocumentType = "weekly"
	DocTypeProject   DocumentType = "project"
	DocTypeMeeting   DocumentType = "meeting"
	DocTypeTemplate  DocumentType = "template"
	DocTypeIndex     DocumentType = "index"
	DocTypeLog       DocumentType = "log"
	DocTypeDatabase  DocumentType = "database"
	DocTypeKanban    DocumentType = "kanban"
	DocTypeTimeline  DocumentType = "timeline"
	DocTypeJournal   DocumentType = "journal"
	DocTypeReference DocumentType = "reference"
	DocTypeGeneric   DocumentType = "generic"
)

// BreadcrumbItem represents one segment in a breadcrumb path.
type BreadcrumbItem struct {
	Label string `json:"label"`
	URI   string `json:"uri,omitempty"`
	Kind  string `json:"kind"` // "workspace", "folder", "document", "heading"
	Icon  string `json:"icon"`
}

// DocumentInfo contains detected metadata about a document.
type DocumentInfo struct {
	Type        DocumentType    `json:"type"`
	Title       string          `json:"title"`
	Breadcrumbs []BreadcrumbItem `json:"breadcrumbs"`
	Parent      string          `json:"parent,omitempty"`
	Tags        []string        `json:"tags,omitempty"`
	Project     string          `json:"project,omitempty"`
	Date        string          `json:"date,omitempty"`
	Author      string          `json:"author,omitempty"`
	Icon        string          `json:"icon"`
}

// DetectDocumentType analyzes a document to determine its type from
// frontmatter, filename patterns, content structure, and path.
func (s *State) DetectDocumentType(uri string) *DocumentInfo {
	text, ok := s.Documents[uri]
	if !ok {
		return &DocumentInfo{Type: DocTypeGeneric, Title: filepath.Base(uri)}
	}

	info := &DocumentInfo{
		Type: DocTypeGeneric,
	}

	// Extract frontmatter metadata
	fm := parseFrontmatterFields(text)

	// Check explicit type in frontmatter
	if t, ok := fm["type"]; ok {
		switch strings.ToLower(t) {
		case "note":
			info.Type = DocTypeNote
		case "daily", "day":
			info.Type = DocTypeDaily
		case "weekly", "week":
			info.Type = DocTypeWeekly
		case "project":
			info.Type = DocTypeProject
		case "meeting":
			info.Type = DocTypeMeeting
		case "template":
			info.Type = DocTypeTemplate
		case "log":
			info.Type = DocTypeLog
		case "database", "db":
			info.Type = DocTypeDatabase
		case "kanban":
			info.Type = DocTypeKanban
		case "timeline":
			info.Type = DocTypeTimeline
		case "journal":
			info.Type = DocTypeJournal
		case "reference", "ref":
			info.Type = DocTypeReference
		case "index":
			info.Type = DocTypeIndex
		}
	}

	// Extract other frontmatter fields
	if title, ok := fm["title"]; ok {
		info.Title = title
	}
	if tags, ok := fm["tags"]; ok {
		for _, t := range strings.Split(tags, ",") {
			t = strings.TrimSpace(strings.Trim(t, "[]\"' "))
			if t != "" {
				info.Tags = append(info.Tags, t)
			}
		}
	}
	if project, ok := fm["project"]; ok {
		info.Project = project
	}
	if date, ok := fm["date"]; ok {
		info.Date = date
	}
	if author, ok := fm["author"]; ok {
		info.Author = author
	}
	if parent, ok := fm["parent"]; ok {
		info.Parent = parent
	}

	// If type not set by frontmatter, detect from path/filename
	if info.Type == DocTypeGeneric {
		info.Type = detectTypeFromPath(uri)
	}

	// If type still generic, detect from content structure
	if info.Type == DocTypeGeneric {
		info.Type = detectTypeFromContent(text)
	}

	// Set title from first heading if not in frontmatter
	if info.Title == "" {
		info.Title = extractFirstHeading(text)
	}
	if info.Title == "" {
		info.Title = strings.TrimSuffix(filepath.Base(strings.TrimPrefix(uri, "file://")), filepath.Ext(uri))
	}

	// Set icon based on type
	info.Icon = iconForType(info.Type)

	// Build breadcrumbs
	info.Breadcrumbs = s.buildBreadcrumbs(uri, info)

	return info
}

// parseFrontmatterFields extracts key-value pairs from YAML frontmatter.
func parseFrontmatterFields(text string) map[string]string {
	fields := make(map[string]string)
	lines := strings.Split(text, "\n")
	if len(lines) < 2 || strings.TrimSpace(lines[0]) != "---" {
		return fields
	}

	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "---" {
			break
		}
		idx := strings.Index(line, ":")
		if idx > 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			fields[strings.ToLower(key)] = val
		}
	}
	return fields
}

// detectTypeFromPath infers document type from its file path.
func detectTypeFromPath(uri string) DocumentType {
	path := strings.ToLower(strings.TrimPrefix(uri, "file://"))
	base := filepath.Base(path)
	dir := filepath.Dir(path)

	// Date-based filenames → daily/journal
	if len(base) >= 10 {
		// YYYY-MM-DD pattern
		prefix := base[:10]
		if len(prefix) == 10 && prefix[4] == '-' && prefix[7] == '-' {
			allDigits := true
			for _, i := range []int{0, 1, 2, 3, 5, 6, 8, 9} {
				if prefix[i] < '0' || prefix[i] > '9' {
					allDigits = false
					break
				}
			}
			if allDigits {
				return DocTypeDaily
			}
		}
	}

	// Special filenames
	nameNoExt := strings.TrimSuffix(base, filepath.Ext(base))
	switch nameNoExt {
	case "index", "readme", "home":
		return DocTypeIndex
	case "changelog", "log":
		return DocTypeLog
	}

	// Directory-based detection
	dirBase := filepath.Base(dir)
	switch dirBase {
	case "daily", "days":
		return DocTypeDaily
	case "weekly", "weeks":
		return DocTypeWeekly
	case "journal", "journals":
		return DocTypeJournal
	case "projects":
		return DocTypeProject
	case "meetings":
		return DocTypeMeeting
	case "templates":
		return DocTypeTemplate
	case "logs":
		return DocTypeLog
	case "references", "refs":
		return DocTypeReference
	}

	return DocTypeGeneric
}

// detectTypeFromContent infers document type from content patterns.
func detectTypeFromContent(text string) DocumentType {
	lines := strings.Split(text, "\n")
	taskCount := 0
	headingCount := 0
	hasKanbanHeadings := false
	todoHeading := false
	inProgressHeading := false
	doneHeading := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			headingCount++
			lower := strings.ToLower(trimmed)
			if strings.Contains(lower, "todo") || strings.Contains(lower, "to do") || strings.Contains(lower, "backlog") {
				todoHeading = true
			}
			if strings.Contains(lower, "in progress") || strings.Contains(lower, "doing") || strings.Contains(lower, "active") {
				inProgressHeading = true
			}
			if strings.Contains(lower, "done") || strings.Contains(lower, "completed") || strings.Contains(lower, "finished") {
				doneHeading = true
			}
		}
		if strings.Contains(trimmed, "- [ ]") || strings.Contains(trimmed, "- [x]") || strings.Contains(trimmed, "- [X]") {
			taskCount++
		}
	}

	if todoHeading && (inProgressHeading || doneHeading) {
		hasKanbanHeadings = true
	}

	if hasKanbanHeadings && taskCount > 3 {
		return DocTypeKanban
	}

	if taskCount > 5 && headingCount <= 2 {
		return DocTypeProject
	}

	return DocTypeGeneric
}

// extractFirstHeading returns the text of the first H1 heading.
func extractFirstHeading(text string) string {
	lines := strings.Split(text, "\n")
	inFM := false
	for i, line := range lines {
		if i == 0 && strings.TrimSpace(line) == "---" {
			inFM = true
			continue
		}
		if inFM {
			if strings.TrimSpace(line) == "---" {
				inFM = false
			}
			continue
		}
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			return strings.TrimSpace(trimmed[2:])
		}
	}
	return ""
}

// buildBreadcrumbs constructs a navigation breadcrumb path for the document.
func (s *State) buildBreadcrumbs(uri string, info *DocumentInfo) []BreadcrumbItem {
	var crumbs []BreadcrumbItem

	path := strings.TrimPrefix(uri, "file://")

	// Find which workspace this belongs to
	// Use directory segments as breadcrumbs
	parts := strings.Split(filepath.Dir(path), string(filepath.Separator))

	// Add workspace root (last meaningful directory)
	if len(parts) > 0 {
		// Find workspace-level dir (look for .down or known workspace indicators)
		wsIdx := 0
		for i := len(parts) - 1; i >= 0; i-- {
			if parts[i] == "notes" || parts[i] == "journal" || parts[i] == "projects" ||
				parts[i] == "templates" || parts[i] == "daily" || parts[i] == "meetings" {
				wsIdx = i
				break
			}
		}
		if wsIdx > 0 && wsIdx < len(parts) {
			// Workspace name is the parent of the special dir
			wsName := parts[wsIdx-1]
			crumbs = append(crumbs, BreadcrumbItem{
				Label: wsName,
				Kind:  "workspace",
				Icon:  "📂",
			})
		}

		// Add intermediate folders after workspace
		for i := wsIdx; i < len(parts); i++ {
			if parts[i] != "" && parts[i] != "." {
				crumbs = append(crumbs, BreadcrumbItem{
					Label: parts[i],
					Kind:  "folder",
					Icon:  "📁",
				})
			}
		}
	}

	// Add the document itself
	crumbs = append(crumbs, BreadcrumbItem{
		Label: info.Title,
		URI:   uri,
		Kind:  "document",
		Icon:  info.Icon,
	})

	return crumbs
}

// iconForType returns an emoji icon for a document type.
func iconForType(dt DocumentType) string {
	switch dt {
	case DocTypeNote:
		return "📝"
	case DocTypeDaily:
		return "📅"
	case DocTypeWeekly:
		return "📆"
	case DocTypeProject:
		return "📋"
	case DocTypeMeeting:
		return "🤝"
	case DocTypeTemplate:
		return "📄"
	case DocTypeIndex:
		return "🏠"
	case DocTypeLog:
		return "📜"
	case DocTypeDatabase:
		return "🗃️"
	case DocTypeKanban:
		return "📊"
	case DocTypeTimeline:
		return "⏳"
	case DocTypeJournal:
		return "📓"
	case DocTypeReference:
		return "📚"
	default:
		return "📄"
	}
}
