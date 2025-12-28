package detector

import (
	"deaiify/internal/parser"
)

// FileScore represents the AI-likeness score for a file
type FileScore struct {
	Path          string
	TotalComments int
	AIComments    int
	Score         float64 // 0.0 to 1.0
	Details       []CommentScore
}

// CommentScore holds detection result for a single comment
type CommentScore struct {
	Comment parser.Comment
	Result  DetectionResult
}

// ScoreFile analyzes all comments in a file and returns an AI-likeness score
func ScoreFile(path string, comments []parser.Comment) FileScore {
	score := FileScore{
		Path:          path,
		TotalComments: len(comments),
		Details:       make([]CommentScore, 0, len(comments)),
	}

	if len(comments) == 0 {
		return score
	}

	var totalScore float64
	for _, comment := range comments {
		result := DetectAIComment(comment)
		score.Details = append(score.Details, CommentScore{
			Comment: comment,
			Result:  result,
		})
		if result.IsAILike {
			score.AIComments++
		}
		totalScore += result.Score
	}

	score.Score = totalScore / float64(len(comments))
	return score
}

// GetAIComments returns only the comments detected as AI-like
func (fs FileScore) GetAIComments() []CommentScore {
	var result []CommentScore
	for _, cs := range fs.Details {
		if cs.Result.IsAILike {
			result = append(result, cs)
		}
	}
	return result
}
