package imagehash

import (
	"fmt"
	"image"

	"github.com/nfnt/resize"
)

// Hasher ImageHashに関するパラメータ群
type Hasher struct {
	width    uint8
	height   uint8
	hashType HashType
}

// Result ImageHashの結果
type Result struct {
	ResizedImg image.Image
	HashVal    []uint8
}

// HashType ImageHashのタイプ
type HashType int

// HashType ImageHashのタイプ
const (
	DHash = 1
	AHash = 2
)

//BuildImageHash ImageHasherを生成する
func BuildImageHash(width uint8, height uint8, hashType HashType) Hasher {
	return Hasher{
		width:    width,
		height:   height,
		hashType: hashType,
	}
}

// CalcHash ハッシュを計算する
func (imageHash *Hasher) CalcHash(img *image.Image) Result {
	resizedImg := resize.Resize(uint(imageHash.width), uint(imageHash.height), *img, resize.Bilinear)
	var result Result
	result.ResizedImg = resizedImg
	switch imageHash.hashType {
	case DHash:
		result.HashVal = calcDHash(&resizedImg)
	case AHash:
		result.HashVal = calcAHash(&resizedImg)
	default:
		panic(fmt.Sprintf("Unknown hash type : %d", imageHash.hashType))
	}
	return result
}

func calcDHash(resizedImg *image.Image) []uint8 {
	width := (*resizedImg).Bounds().Max.X
	height := (*resizedImg).Bounds().Max.Y
	hashArr := make([]uint8, width*height)
	var prevPixVal uint32 = 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := (*resizedImg).At(x, y).RGBA()
			pixVal := (r + g + b) / 3
			if prevPixVal < pixVal {
				hashArr[height*y+x] = 1
			} else {
				hashArr[height*y+x] = 0
			}
			prevPixVal = pixVal
		}
	}
	return hashArr
}

func calcAHash(resizedImg *image.Image) []uint8 {
	// 未実装
	width := (*resizedImg).Bounds().Max.X
	height := (*resizedImg).Bounds().Max.Y
	hashArr := make([]uint8, width*height)
	avgVal := calcAveragePixVal(resizedImg)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := (*resizedImg).At(x, y).RGBA()
			pixVal := (r + g + b) / 3
			if uint32(avgVal) < pixVal {
				hashArr[height*y+x] = 1
			} else {
				hashArr[height*y+x] = 0
			}
		}
	}
	return hashArr
}

func calcAveragePixVal(resizedImg *image.Image) uint8 {
	width := (*resizedImg).Bounds().Max.X
	height := (*resizedImg).Bounds().Max.Y
	totalVal := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := (*resizedImg).At(x, y).RGBA()
			pixVal := (r + g + b) / 3
			totalVal += int(pixVal)
		}
	}
	return uint8(totalVal / (width * height))
}

// CalcSimilarity ハッシュ同士の類似度を計算する
func (result1 *Result) CalcSimilarity(result2 *Result) (float32, error) {
	if len(result1.HashVal) != len(result2.HashVal) {
		return 0, fmt.Errorf("hash length is not same. %d and %d", len(result1.HashVal), len(result2.HashVal))
	}
	sameCount := 0
	for i := 0; i < len(result1.HashVal); i++ {
		if result1.HashVal[i] == result2.HashVal[i] {
			sameCount++
		}
	}
	return (float32(sameCount) / float32(len(result1.HashVal))), nil
}
