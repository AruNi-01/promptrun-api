package cron

func StarCronJobs() {
	syncAverageRatingJob()
	calculatePromptScoreJob()
}
