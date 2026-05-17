package ui

import "image"

var (
	// Main Window
	ScreenWidth  = 1280
	ScreenHeight = 768

	// Central Widget
	WidgetWidth  = 900
	WidgetHeight = 600
	WidgetX      = (ScreenWidth - WidgetWidth) / 2
	WidgetY      = (ScreenHeight - WidgetHeight) / 2 + 10
	WidgetRect   = image.Rect(WidgetX, WidgetY, WidgetX+WidgetWidth, WidgetY+WidgetHeight)

	// Tabs
	TabHeight  = 40
	TabRect    = image.Rect(WidgetX, WidgetY, WidgetX+WidgetWidth, WidgetY+TabHeight)
	TabWidth   = WidgetWidth / 3
	Tab1Rect   = image.Rect(WidgetX, WidgetY, WidgetX+TabWidth, WidgetY+TabHeight)
	Tab2Rect   = image.Rect(WidgetX+TabWidth, WidgetY, WidgetX+TabWidth*2, WidgetY+TabHeight)
	Tab3Rect   = image.Rect(WidgetX+TabWidth*2, WidgetY, WidgetX+WidgetWidth, WidgetY+TabHeight)

	// Content Area (below tabs)
	ContentRect = image.Rect(WidgetX, WidgetY+TabHeight, WidgetX+WidgetWidth, WidgetY+WidgetHeight)

	// HUD (Top of content)
	MetricsHUDRect = image.Rect(WidgetX+10, WidgetY+TabHeight+10, WidgetX+WidgetWidth-10, WidgetY+TabHeight+110)

	// Interaction Area (Bottom of content)
	// Hardware/Upgrade list
	ListRect = image.Rect(WidgetX+10, WidgetY+TabHeight+120, WidgetX+WidgetWidth-10, WidgetY+WidgetHeight-10)

	// Clicker Region (Specific to a tab or always visible? Let's put it in a "DASHBOARD" tab or SYSTEM?)
	// Let's redefine tabs: "DASHBOARD", "HARDWARE", "UPGRADES"
	ClickerRegion = image.Rect(WidgetX+WidgetWidth/2-100, WidgetY+TabHeight+250, WidgetX+WidgetWidth/2+100, WidgetY+TabHeight+330)

	// Reboot Button (in SYSTEM tab)
	RebootBtnRect = image.Rect(WidgetX+10, WidgetY+WidgetHeight-50, WidgetX+WidgetWidth-10, WidgetY+WidgetHeight-10)
)
