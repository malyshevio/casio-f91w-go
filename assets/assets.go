package casio_f91w_assets


import (
	"embed"
)

//go:embed beepNoTrailingDelay.wav
//go:embed beepBeepNoTrailingDelay.wav
//go:embed beepBeepWithTrailingDelay.wav
var Assets embed.FS
