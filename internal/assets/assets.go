// Package assets contains embedded assets
package assets

import "embed"

// LicenseFS is an [embed.FS] with every embedded License files.
//
//go:embed licenses/*
var LicenseFS embed.FS

//go:embed plugins/*
var PluginFS embed.FS
