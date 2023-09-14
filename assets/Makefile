default: beepBeepNoTrailingDelay.pcm beepBeepWithTrailingDelay.pcm beepNoTrailingDelay.pcm

%.pcm: %.wav
	ffmpeg -loglevel error -i $< -f s16le -acodec pcm_s16le $@

