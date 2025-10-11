package leetcode_graphql

import (
	"bytes"
	"encoding/json"
	"io"
	"main/routes/utils"
	"net/http"
	"strings"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

const LEETCODE_GRAPHQL_URL = "https://leetcode.com/graphql"

func (h *Handler) GetOnlineUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		utils.WriteInternalServerErrorResponse(w, "Username parameter required")
		return
	}

	query := `
		query getUserProfile($username: String!) {
			matchedUser(username: $username) {
				username
				submitStats: submitStatsGlobal {
					acSubmissionNum {
						difficulty
						count
						submissions
					}
				}
			}
		}
	`

	requestBody := map[string]interface{}{
		"query": query,
		"variables": map[string]string{
			"username": username,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		utils.WriteInternalServerErrorResponse(w, "Failed to create request")
		return
	}

	resp, err := http.Post(LEETCODE_GRAPHQL_URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		utils.WriteInternalServerErrorResponse(w, "Request failed: "+err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.WriteInternalServerErrorResponse(w, "Failed to read response")
		return
	}

	if resp.StatusCode != http.StatusOK {
		utils.WriteInternalServerErrorResponse(w, "Error occurred: "+string(body))
		return
	}

	var result GetOnlineUsernameResponse
	if err := json.Unmarshal(body, &result); err != nil {
		utils.WriteInternalServerErrorResponse(w, "Failed to parse response")
		return
	}

	utils.WriteJSONResponse(w, result.Data.MatchedUser, http.StatusOK)
}

func (h *Handler) GetDifficulty(w http.ResponseWriter, r *http.Request) {
	link := r.URL.Query().Get("link")
	if link == "" {
		utils.WriteInternalServerErrorResponse(w, "Link parameter required")
		return
	}

	parts := strings.Split(link, "/problems/")
	if len(parts) < 2 {
		utils.WriteInternalServerErrorResponse(w, "Invalid problem link format")
		return
	}

	titleSlug := strings.Split(parts[1], "/")[0]

	query := `
		query getProblemDetails($titleSlug: String!) {
			question(titleSlug: $titleSlug) {
				title
				difficulty
			}
		}
	`

	requestBody := map[string]interface{}{
		"query": query,
		"variables": map[string]string{
			"titleSlug": titleSlug,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		utils.WriteInternalServerErrorResponse(w, "Failed to create request")
		return
	}

	resp, err := http.Post(LEETCODE_GRAPHQL_URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		utils.WriteInternalServerErrorResponse(w, "Request failed: "+err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.WriteInternalServerErrorResponse(w, "Failed to read response")
		return
	}

	if resp.StatusCode != http.StatusOK {
		utils.WriteInternalServerErrorResponse(w, "Error occurred: "+string(body))
		return
	}

	var result GetDifficultyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		utils.WriteInternalServerErrorResponse(w, "Failed to parse response")
		return
	}

	utils.WriteJSONResponse(w, result.Data.Question, http.StatusOK)
}

func (h *Handler) ExtractProblem(w http.ResponseWriter, r *http.Request) {
	link := r.URL.Query().Get("link")
	if link == "" {
		utils.WriteInternalServerErrorResponse(w, "Link parameter required")
		return
	}

	if !strings.Contains(link, "/submissions") || !strings.Contains(link, "/problems/") {
		utils.WriteJSONResponse(w, ExtractProblemResponse{Problem: ""}, http.StatusOK)
		return
	}

	parts := strings.Split(link, "/problems/")
	if len(parts) < 2 {
		utils.WriteJSONResponse(w, ExtractProblemResponse{Problem: ""}, http.StatusOK)
		return
	}

	problemSlug := strings.Split(parts[1], "/submissions")[0]
	problem := strings.ReplaceAll(problemSlug, "-", " ")

	utils.WriteJSONResponse(w, ExtractProblemResponse{Problem: problem}, http.StatusOK)
}

