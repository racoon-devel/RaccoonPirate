package selector

type Settings struct {
	MinSeasonSizeMB     int64
	MaxSeasonSizeMB     int64
	MinSeedersThreshold int64
	Voice               string
	VoiceList           Voices
	QualityPrior        []string
	AudioFormats        []string
}
