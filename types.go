package main

type Options struct {
	Workdir        string
	TargetBucket   string
	TargetPath     string
	Repository     string
	RepositoryList []string
	GithubOptions  GithubCliOptions
}

type GithubCliOptions struct {
	Mirror bool
}

type GitCommand struct {
	options Options
}

type TarCommand struct {
	options Options
}

type S3Command struct {
	options Options
}

type CleanCommand struct {
	options Options
}
