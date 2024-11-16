package main

import (
	"html/template"
	"time"
)

func setTemplateHelpers() template.FuncMap {
	return template.FuncMap{
		"unescapeHTML":    unescapeHTML,
		"formatDate":      formatDate,
		"formatDateLong":  formatDateLong,
		"formatDateShort": formatDateShort,
	}
}

func unescapeHTML(s string) template.HTML {
	return template.HTML(s)
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func formatDateLong(t time.Time) string {
	return t.Format("2 Jan, 2006")
}

func formatDateShort(t time.Time) string {
	return t.Format("02-01-2006")
}
