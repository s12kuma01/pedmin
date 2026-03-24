// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// BuildAvatarGallery builds the avatar display UI with server and global avatars.
func BuildAvatarGallery(user discord.User, member *discord.ResolvedMember, guildID *snowflake.ID) discord.ContainerComponent {
	cdnOpts := []discord.CDNOpt{
		discord.WithFormat(discord.FileFormatPNG),
		discord.WithSize(1024),
	}

	displayName := user.EffectiveName()
	globalURL := user.EffectiveAvatarURL(cdnOpts...)

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay(fmt.Sprintf("### %s のアバター", displayName)),
	}

	var serverURL *string
	if member != nil && guildID != nil {
		// Ensure GuildID is set for MemberAvatar CDN path
		m := member.Member
		m.GuildID = *guildID
		serverURL = m.AvatarURL(cdnOpts...)
	}

	if serverURL != nil && *serverURL != globalURL {
		components = append(components,
			discord.NewTextDisplay("**サーバーアバター**"),
			discord.NewMediaGallery(discord.MediaGalleryItem{
				Media: discord.UnfurledMediaItem{URL: *serverURL},
			}),
			discord.NewLargeSeparator(),
			discord.NewTextDisplay("**グローバルアバター**"),
			discord.NewMediaGallery(discord.MediaGalleryItem{
				Media: discord.UnfurledMediaItem{URL: globalURL},
			}),
		)
	} else {
		components = append(components,
			discord.NewMediaGallery(discord.MediaGalleryItem{
				Media: discord.UnfurledMediaItem{URL: globalURL},
			}),
		)
	}

	return discord.NewContainer(components...)
}
