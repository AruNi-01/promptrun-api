package cron

func JobsInit() {
	syncAverageRatingJob()
	refreshPromptScoreJob()
}
