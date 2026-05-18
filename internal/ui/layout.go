package ui

import "image"

var (
	// Main Window
	ScreenWidth  = 1280
	ScreenHeight = 768

	// Central Widget
	WidgetWidth  = 1000
	WidgetHeight = 600
	WidgetX      = (ScreenWidth - WidgetWidth) / 2
	WidgetY      = (ScreenHeight - WidgetHeight) / 2 + 50
	WidgetRect   = image.Rect(WidgetX, WidgetY, WidgetX+WidgetWidth, WidgetY+WidgetHeight)

	// HUD (Top of content)
	MetricsHUDRect = image.Rect(WidgetX+20, WidgetY+20, WidgetX+WidgetWidth-20, WidgetY+140)

	// Clicker Region (Middle-Left)
	ClickerRegion = image.Rect(WidgetX+20, WidgetY+160, WidgetX+320, WidgetY+300)

	// Hardware/Upgrade List (Middle-Right)
	ListRect = image.Rect(WidgetX+340, WidgetY+160, WidgetX+WidgetWidth-20, WidgetY+500)

	// Log Rect (Bottom-Left)
	LogRect = image.Rect(WidgetX+20, WidgetY+320, WidgetX+320, WidgetY+500)

	// Packet Intercept (Bottom-Centerish, left of Reboot)
	// Log ends at WidgetX+320, Reboot starts at WidgetX+340. 
	// To be "left of prestige button" and "under log" is tight.
	// Let's place it at WidgetX+20, WidgetY+520 (under log) but sized to fit left of Reboot.
	PacketRect = image.Rect(WidgetX+20, WidgetY+520, WidgetX+320, WidgetY+580)

	// Reboot Button (Bottom-Right)
	RebootBtnRect = image.Rect(WidgetX+340, WidgetY+520, WidgetX+WidgetWidth-20, WidgetY+580)
)
