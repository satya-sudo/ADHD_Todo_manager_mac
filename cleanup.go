package main

import "time"

func cleanupOldTasks() {

	now := time.Now().Unix()

	var fresh []Task

	for _, t := range tasks {

		if now-t.CreatedAt < 86400 || t.State != Done {
			fresh = append(fresh, t)
		}

	}

	tasks = fresh

	saveTasks()
}

func startCleanupLoop() {

	go func() {

		for {

			time.Sleep(1 * time.Hour)

			cleanupOldTasks()

		}

	}()

}
