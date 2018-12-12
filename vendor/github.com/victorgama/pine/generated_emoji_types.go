package pine

var typeEmoji = map[msgType]string{
	Info:      "💬",
	Success:   "✅",
	Warn:      "⚠️",
	Error:     "🚨",
	Timing:    "⏱",
	WTF:       "👻",
	Finish:    "🏁",
	Terminate: "⛔️",
	Spawn:     "✨",
	Disk:      "💾",
}

var typeMap = map[msgType]string{
	Spawn:     "Spawn",
	Disk:      "Disk",
	WTF:       "WTF",
	Finish:    "Finish",
	Terminate: "Terminate",
	Error:     "Error",
	Timing:    "Timing",
	Info:      "Info",
	Success:   "Success",
	Warn:      "Warn",
}
