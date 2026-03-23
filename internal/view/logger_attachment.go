// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/ui"
)

// LoggerBuildAttachmentComponents builds container sub-components for a list of attachments.
func LoggerBuildAttachmentComponents(attachments []discord.Attachment) []discord.ContainerSubComponent {
	var images []discord.MediaGalleryItem
	var files []string

	for _, a := range attachments {
		if a.ContentType != nil && strings.HasPrefix(*a.ContentType, "image/") {
			images = append(images, discord.MediaGalleryItem{
				Media: discord.UnfurledMediaItem{URL: a.URL},
			})
		} else {
			size := ui.FormatBytes(uint64(a.Size))
			files = append(files, fmt.Sprintf("📎 %s (%s)", a.Filename, size))
		}
	}

	var components []discord.ContainerSubComponent
	if len(images) > 0 {
		components = append(components, discord.NewMediaGallery(images...))
	}
	if len(files) > 0 {
		components = append(components, discord.NewTextDisplay(strings.Join(files, "\n")))
	}
	return components
}

// LoggerDiffAttachments computes the removed and added attachments between old and new.
func LoggerDiffAttachments(old, new []discord.Attachment) (removed, added []discord.Attachment) {
	oldIDs := make(map[snowflake.ID]discord.Attachment, len(old))
	for _, a := range old {
		oldIDs[a.ID] = a
	}
	newIDs := make(map[snowflake.ID]struct{}, len(new))
	for _, a := range new {
		newIDs[a.ID] = struct{}{}
		if _, exists := oldIDs[a.ID]; !exists {
			added = append(added, a)
		}
	}
	for _, a := range old {
		if _, exists := newIDs[a.ID]; !exists {
			removed = append(removed, a)
		}
	}
	return
}

// LoggerAttachmentsEqual checks if two attachment slices have the same IDs in order.
func LoggerAttachmentsEqual(old, new []discord.Attachment) bool {
	if len(old) != len(new) {
		return false
	}
	for i := range old {
		if old[i].ID != new[i].ID {
			return false
		}
	}
	return true
}
