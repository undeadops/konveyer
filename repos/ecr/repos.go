package ecr

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/undeadops/konveyer/common"
)

// ContainerImages - Information about Available Containers
type ContainerImages struct {
	ImageDigest      *string    `locationName:"imageDigest" type:"string"`
	ImagePushedAt    *time.Time `locationName:"imagePushedAt" type:"timestamp" timestampFormat:"unix"`
	ImageSizeInBytes *int64     `locationName:"imageSizeInBytes" type:"long"`
	ImageTags        []string   `locationName:"imageTags" type:"list"`
	RepositoryName   *string    `locationName:"repositoryName" min:"2" type:"string"`
}

// Repository - AWS ECR Repositories
type Repository struct {
	CreatedAt      *time.Time `locationName:"createdAt" type:"timestamp" timestampFormat:"unix"`
	RepositoryArn  *string    `locationName:"repositoryArn" type:"string"`
	RepositoryName *string    `locationName:"repositoryName" min:"2" type:"string"`
	RepositoryUri  *string    `locationName:"repositoryUri" type:"string"`
}

// FetchContainerRepos - Return List of Repos of Containers
func FetchContainerRepos() []Repository {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Panicf("panic: failed to load config, " + err.Error())
	}
	cfg.Region = "us-east-2"

	svc := ecr.New(cfg)
	input := &ecr.DescribeRepositoriesInput{}
	req := svc.DescribeRepositoriesRequest(input)
	result, err := req.Send()
	if err != nil {
		log.Panicf("panic: failed to query aws DescribeRepositoriesRequest," + err.Error())
	}

	repos := []Repository{}

	for _, repo := range result.Repositories {
		r := Repository{
			CreatedAt:      repo.CreatedAt,
			RepositoryArn:  repo.RepositoryArn,
			RepositoryName: repo.RepositoryName,
			RepositoryUri:  repo.RepositoryUri,
		}
		repos = append(repos, r)
	}

	return repos
}

// CreateContainerRepo - Create Container Repo with project creation
func CreateContainerRepo(name string) (Repository, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Panicf("panic: Failed to load config, " + err.Error())
	}
	cfg.Region = "us-east-2"

	svc := ecr.New(cfg)

	input := &ecr.CreateRepositoryInput{
		RepositoryName: &name,
	}
	req := svc.CreateRepositoryRequest(input)
	result, err := req.Send()
	if err != nil {
		log.Panicf("panic: Failed to Create ECR Repo, " + err.Error())
		return Repository{}, err
	}

	r := Repository{
		CreatedAt:      result.Repository.CreatedAt,
		RepositoryArn:  result.Repository.RepositoryArn,
		RepositoryName: result.Repository.RepositoryName,
		RepositoryUri:  result.Repository.RepositoryUri,
	}

	return r, nil
}

// FetchContainerVersions - Return List of Versions of a speicic Container Repo
func FetchContainerVersions(repo string) ([]ContainerImages, error) {
	r := common.GetRuntime()

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Panicf("panic: Failed to load config, " + err.Error())
	}
	cfg.Region = r.Region

	svc := ecr.New(cfg)

	input := &ecr.DescribeImagesInput{
		RepositoryName: &repo,
	}
	req := svc.DescribeImagesRequest(input)
	result, err := req.Send()
	if err != nil {
		log.Panicf("panic: Failed to List Containers for Project, " + err.Error())
		return nil, err
	}

	images := []ContainerImages{}
	for _, r := range result.ImageDetails {
		c := ContainerImages{
			ImageDigest:      r.ImageDigest,
			ImagePushedAt:    r.ImagePushedAt,
			ImageSizeInBytes: r.ImageSizeInBytes,
			ImageTags:        r.ImageTags,
			RepositoryName:   r.RepositoryName,
		}
		images = append(images, c)
	}

	return images, nil
}

// VerifyContainerImageVersion - Check if valid Container Image Tag
func VerifyContainerImageVersion(repo Repository, tag string) bool {
	r := common.GetRuntime()

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Panicf("panic: Failed to load config, " + err.Error())
		return false
	}
	cfg.Region = r.Region

	svc := ecr.New(cfg)

	input := &ecr.DescribeImagesInput{
		RepositoryName: repo.RepositoryName,
	}
	req := svc.DescribeImagesRequest(input)
	result, err := req.Send()
	if err != nil {
		log.Panicf("panic: Failed to List Containers for Repo, " + err.Error())
		return false
	}

	// Loop over Images and Image Tags, For Image with tag
	for _, i := range result.ImageDetails {
		for _, imageTags := range i.ImageTags {
			if imageTags == tag {
				return true
			}
		}
	}
	return false
}
