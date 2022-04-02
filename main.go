package main

import (
	"bytes"
	"duplicate_image_detector/imagehash"
	"duplicate_image_detector/numberutil"
	"duplicate_image_detector/resultwriter"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var (
	dirPath     string
	verbose     bool
	threshold   float64
	outFileName string
)

type hashWithFileInfo struct {
	file   fs.FileInfo
	result imagehash.Result
}

func main() {
	flag.StringVar(&dirPath, "dir", "", "directory path")
	flag.BoolVar(&verbose, "verbose", false, "directory path")
	flag.Float64Var(&threshold, "threshold", 1.0, "similarity threshold")
	flag.StringVar(&outFileName, "out", "out", "out file name without extension")
	flag.Parse()
	if dirPath == "" {
		fmt.Println("directory path is empty")
		os.Exit(1)
	}

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		fmt.Printf("%s is not exist\n", dirPath)
		os.Exit(1)
	}

	hashList, err := calcAllHash()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	pairInfoList, err := calcSimlarities(hashList)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = resultwriter.WriteResultToHTML(pairInfoList, outFileName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = resultwriter.WriteResultToCSV(pairInfoList, outFileName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func calcAllHash() ([]hashWithFileInfo, error) {
	fmt.Println("Preprocessing...")

	hasher := imagehash.BuildImageHash(8, 8, imagehash.DHash)
	result := []hashWithFileInfo{}
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	totalCount := len(files)
	processingCount := 0
	validExtentions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	for _, file := range files {
		_, validExtension := validExtentions[filepath.Ext(file.Name())]
		if file.IsDir() || !validExtension {
			totalCount--
			fmt.Printf("\r %d/%d  ", processingCount, totalCount)
			continue
		}

		image, err := readImage(filepath.Join(dirPath, file.Name()))
		if err != nil {
			totalCount--
			fmt.Printf("\r %d/%d  ", processingCount, totalCount)
			continue
		}
		processingCount++
		result = append(result, hashWithFileInfo{
			file, hasher.CalcHash(&image),
		})
		fmt.Printf("\r %d/%d  ", processingCount, totalCount)
	}

	fmt.Println("\nPreprocessing finished!!")
	return result, nil
}

func calcSimlarities(fileHashList []hashWithFileInfo) ([]resultwriter.SimilarImageGroupInfo, error) {
	fmt.Println("\nCalculation start!")
	start := time.Now()
	result := []resultwriter.SimilarImageGroupInfo{}
	detectedImgIndexes := map[int]bool{}
	totalCount, _ := numberutil.CombinationCount(uint32(len(fileHashList)), 2)
	processingCount := 0
	for i, originHash := range fileHashList {
		if _, detected := detectedImgIndexes[i]; detected {
			processingCount += len(fileHashList) - i - 1
			continue
		}
		similarImages := []resultwriter.SimilarImageInfo{}
		for j := i + 1; j < len(fileHashList); j++ {
			processingCount++
			targetHash := fileHashList[j]
			sim, err := originHash.result.CalcSimilarity(&targetHash.result)
			if err != nil {
				panic(err.Error())
			}

			if sim >= float32(threshold) {
				similarImages = append(similarImages, resultwriter.SimilarImageInfo{
					Similarity:       sim,
					AbsoluteFilePath: filepath.Join(dirPath, targetHash.file.Name()),
					File:             targetHash.file,
				})
				detectedImgIndexes[j] = true
				fmt.Printf("\n[similar image pair] %s : %s - similarity : %f\n", originHash.file.Name(), fileHashList[j].file.Name(), sim)
			}
			fmt.Printf("\r %d/%d", processingCount, totalCount)
		}
		fmt.Printf("\r %d/%d", processingCount, totalCount)

		if len(similarImages) > 0 {
			result = append(result, resultwriter.SimilarImageGroupInfo{
				OriginFile:       originHash.file,
				AbsoluteFilePath: filepath.Join(dirPath, originHash.file.Name()),
				SimilarImages:    similarImages,
			})
		}
	}
	elapsedTime := time.Now().Sub(start)
	fmt.Printf("\n\nCalculation finished!\nDetected count: %d\n", len(detectedImgIndexes))
	fmt.Printf("Elapsed time: %fs\n", elapsedTime.Seconds())
	return result, nil
}

func readImage(filePath string) (image.Image, error) {
	fileData, err := os.Open(filePath)
	defer fileData.Close()
	if err != nil {
		if verbose {
			fmt.Println("File did not open ", err.Error())
		}
		return nil, err
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, fileData)
	img, _, err := image.Decode(buf)
	if err != nil {
		if verbose {
			fmt.Println("Image decode error ", err.Error(), filePath)
		}
		return nil, err
	}
	return img, nil
}
