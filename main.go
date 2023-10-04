package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// Global variable
type patchMeta struct {
	TagName      string `json:"tagName"`
	URL          string `json:"url"`
	Source       string `json:"source"`
	FrameworkUrl string `json:"frameworkUrl"`
}

type mihoyoApiStruct struct {
	PathData patchMeta `json:"data"`
	Msg      string    `json:"message"`
	Code     int       `json:"retcode"`
}

var (
	GitHubMeta patchMeta
	GitLabMeta patchMeta
	GitHubResponse mihoyoApiStruct
	GitLabResponse mihoyoApiStruct
)

func updateGitHubMeta() {
	for {
		apiURL := "https://api.github.com/repos/HolographicHat/YaeAchievement/releases/latest"
		req, err := http.NewRequest("GET", apiURL, nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return
		}

		// Create Dict object
		var result map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&result); err != nil {
			fmt.Println("JSON Error", err)
			return
		}

		tagName := result["tag_name"].(string)
		assetUrl := result["assets"].([]interface{})[0].(map[string]interface{})["browser_download_url"].(string)

		GitHubMeta.TagName = tagName
		GitHubMeta.URL = assetUrl
		GitHubMeta.Source = "github"

		GitHubResponse.PathData = GitHubMeta
		GitHubResponse.Msg = "success"
		GitHubResponse.Code = 0
		fmt.Println("Successfully build GitHub release cache")

		time.Sleep(2 * time.Minute)
	}
}

func updateGitLabMeta() {
	for {
		apiURL := "https://jihulab.com/api/v4/projects/DGP-Studio%2FYaeAchievement/releases/permalink/latest"
		req, err := http.NewRequest("GET", apiURL, nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return
		}

		// Create Dict object
		var result map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&result); err != nil {
			fmt.Println("JSON Error", err)
			return
		}

		GitLabMeta.TagName = result["tag_name"].(string)
		if assets, ok := result["assets"].(map[string]interface{}); ok {
			if sources, ok := assets["links"].([]interface{}); ok {
				for _, source := range sources {
					if sourceMap, ok := source.(map[string]interface{}); ok {
						if linkType, ok := sourceMap["link_type"].(string); ok && linkType == "package" {
							if name, ok := sourceMap["name"].(string); ok && name == "YaeAchievement.exe" {
								GitLabMeta.URL, _ = sourceMap["direct_asset_url"].(string)
							}
						}
						if name, ok := sourceMap["name"].(string); ok && strings.Contains(name, "Desktop Runtime") {
							FUrl, _ := sourceMap["url"].(string)
							GitLabMeta.FrameworkUrl = FUrl
							GitHubMeta.FrameworkUrl = strings.Replace(FUrl, "zh-cn", "en-us", 1)
							GitHubResponse.PathData = GitHubMeta
						}
					}
				}
			}
		}
		GitLabMeta.Source = "jihulab"

		GitLabResponse.PathData = GitLabMeta
		GitLabResponse.Msg = "success"
		GitLabResponse.Code = 0
		fmt.Println("Successfully build GitLab release cache")

		time.Sleep(2 * time.Minute)
	}
}

func main() {
	go updateGitHubMeta()
	go updateGitLabMeta()

	r := gin.Default()

	r.GET("/global", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, GitHubResponse)
	})

	r.GET("/cn", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, GitLabResponse)
	})

	if err := r.Run(":8080"); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
