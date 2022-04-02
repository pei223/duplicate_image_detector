package resultwriter

import "io/fs"

// SimilarImageGroupInfo 類似画像グループの情報
type SimilarImageGroupInfo struct {
	OriginFile       fs.FileInfo
	AbsoluteFilePath string
	SimilarImages    []SimilarImageInfo
}

// SimilarImageInfo 類似画像情報
type SimilarImageInfo struct {
	File             fs.FileInfo
	AbsoluteFilePath string
	Similarity       float32
}
