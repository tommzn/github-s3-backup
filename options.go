package main

import "strings"

func (options Options) RepositoryName() string {

	if len(options.Repository) == 0 {
		return ""
	}

	repository := strings.ReplaceAll(options.Repository, "git@github.com:", "")
	repository = strings.ReplaceAll(repository, ".git", "")
	return strings.ReplaceAll(repository, "/", ".")
}

func (options Options) RepositoryPathWithWorkDir() string {

	repoName := options.RepositoryName()
	if len(options.Workdir) != 0 {
		return appendPathSeparator(options.Workdir) + repoName
	}
	return repoName
}

func (options Options) ArchiveTarget() string {

	targetPath := options.TargetBucket
	if !strings.HasPrefix(targetPath, "s3://") {
		targetPath = "s3://" + targetPath
	}
	targetPath = appendPathSeparator(targetPath)

	if len(options.TargetPath) != 0 {
		targetPath = targetPath + options.TargetPath
	}
	targetPath = appendPathSeparator(targetPath)
	return targetPath + options.RepositoryName() + ".tar.gz"
}

func appendPathSeparator(path string) string {
	if !strings.HasSuffix(path, "/") {
		return path + "/"
	}
	return path
}
