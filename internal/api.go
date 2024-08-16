package internal

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"

	"net/http"

	"github.com/google/go-github/v63/github"
)

// createGitHubClient membuat klien GitHub dengan konfigurasi HTTP dan timeout
func createGitHubClient() *github.Client {
	httpClient := &http.Client{
		Timeout: 10 * time.Second, // Set timeout untuk permintaan HTTP
	}
	return github.NewClient(httpClient)
}

// GetLatestRelease mengambil versi rilis terbaru dari GitHub untuk repositori yang diberikan
func GetLatestRelease(repoOwner, repoName, releaseType string) (string, error) {
	client := createGitHubClient()

	// Membuat context dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Ambil semua rilis untuk repositori yang diberikan
	releases, _, err := client.Repositories.ListReleases(ctx, repoOwner, repoName, nil)
	if err != nil {
		return "", fmt.Errorf("error getting releases for %s/%s: %w", repoOwner, repoName, err)
	}

	if len(releases) == 0 {
		return "", fmt.Errorf("no releases found for %s/%s", repoOwner, repoName)
	}

	// Menyaring rilis sesuai dengan tipe yang diinginkan
	var filteredReleases []*github.RepositoryRelease
	for _, release := range releases {
		if release.GetPublishedAt().IsZero() {
			continue
		}
		if releaseType == "stable" && !release.GetPrerelease() {
			filteredReleases = append(filteredReleases, release)
		} else if releaseType == "latest" {
			filteredReleases = append(filteredReleases, release)
		}
	}

	if len(filteredReleases) == 0 {
		return "", fmt.Errorf("no %s release found for %s/%s", releaseType, repoOwner, repoName)
	}

	// Mengurutkan rilis berdasarkan waktu publikasi
	sort.Slice(filteredReleases, func(i, j int) bool {
		return filteredReleases[i].GetPublishedAt().Time.After(filteredReleases[j].GetPublishedAt().Time)
	})

	// Mendapatkan rilis terbaru
	latestRelease := filteredReleases[0]
	version := strings.TrimPrefix(*latestRelease.TagName, "v")
	return version, nil
}

const (
    caddyFormat   = "caddy_%s_%s_%s.tar.gz"
    singBoxFormat = "sing-box-%s-%s-%s.tar.gz"
)

// BuildDownloadURL membangun URL unduhan berdasarkan versi, OS, dan arsitektur
func BuildDownloadURL(repoOwner, repoName, version string) (string, error) {
	if repoOwner == "" || repoName == "" || version == "" {
		return "", fmt.Errorf("invalid input: all parameters (repoOwner, repoName, version, os, arch) must be non-empty")
	}

	var fileNameFormat string
	switch repoName {
	case "caddy":
		fileNameFormat = caddyFormat
	case "sing-box":
		fileNameFormat = singBoxFormat
	default:
		return "", fmt.Errorf("unsupported repoName: %s", repoName)
	}

	fileName := fmt.Sprintf(fileNameFormat, version, runtime.GOOS, runtime.GOARCH)

	// Membuat URL unduhan
	downloadURL := fmt.Sprintf("https://github.com/%s/%s/releases/download/v%s/%s", repoOwner, repoName, version, fileName)

	return downloadURL, nil
}
