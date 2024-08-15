package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v53/github"
)

// GetLatestRelease mengambil versi rilis terbaru dari GitHub untuk repositori yang diberikan
func GetLatestRelease(repoOwner, repoName, releaseType string) (string, error) {
	// Membuat client GitHub dengan autentikasi kosong (untuk repositori publik)
	client := github.NewClient(nil)

	// Ambil semua rilis untuk repositori yang diberikan
	releases, _, err := client.Repositories.ListReleases(context.Background(), repoOwner, repoName, nil)
	if err != nil {
		return "", fmt.Errorf("error getting releases for %s/%s: %w", repoOwner, repoName, err)
	}

	if len(releases) == 0 {
		return "", fmt.Errorf("no releases found for %s/%s", repoOwner, repoName)
	}

	// Menentukan versi rilis terbaru berdasarkan jenis rilis
	var latestRelease *github.RepositoryRelease
	for _, release := range releases {
		releaseTime := release.GetPublishedAt().Time // Konversi ke time.Time
		if releaseType == "stable" && !release.GetPrerelease() {
			if latestRelease == nil || releaseTime.After(latestRelease.GetPublishedAt().Time) {
				latestRelease = release
			}
		} else if releaseType == "latest" {
			if latestRelease == nil || releaseTime.After(latestRelease.GetPublishedAt().Time) {
				latestRelease = release
			}
		}
	}

	if latestRelease == nil {
		return "", fmt.Errorf("no %s release found for %s/%s", releaseType, repoOwner, repoName)
	}
	var version = strings.TrimPrefix(*latestRelease.TagName, "v")
	return version, nil
}

// BuildDownloadURL membangun URL unduhan berdasarkan versi, OS, dan arsitektur
func BuildDownloadURL(repoOwner, repoName, version, os, arch string) string {
	baseURL := "https://github.com"
	fileName := fmt.Sprintf("%s-%s-%s-%s.tar.gz", repoName, version, os, arch)
	return fmt.Sprintf("%s/%s/%s/releases/download/v%s/%s", baseURL, repoOwner, repoName, version, fileName)
}
