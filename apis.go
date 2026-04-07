package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/apiFunctions"
	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/constants"
	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/extractionConfig"
	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/sql"

	"github.com/gin-gonic/gin"
)

func main() {
	dbObj, err := sql.NewDatabase(constants.SqlDBName)
	if err != nil {
		log.Fatal(err)
	}
	dbObj.DB.Close()

	r := setupAPIS()
	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}

func setupAPIS() *gin.Engine {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	r.POST(constants.ApiPostImportPath, apiPostImport)
	r.GET(constants.ApiGetRecordsByIdPath, apiGetRecordsById)
	r.GET(constants.ApiGetRecordsPath, apiGetRecords)
	r.POST(constants.ApiPostTransformPath, apiPostTransform)
	r.GET(constants.ApiGetAnalytics, apiGetAnalytics)

	return r
}

// adds items from JSONL received in the request body.
func apiPostImport(c *gin.Context) {
	scanner := bufio.NewScanner(c.Request.Body)
	defer c.Request.Body.Close()

	importReturn := ImportReturn{
		TotalLinesProcessed:         0,
		RecordsImportedSuccessfully: 0,
		ValidationErrors:            make([]string, 0),
		DataQualityWarnings:         make(map[string][]string),
		Statistics: Statistics{
			RecordsByType:  make(map[string]int),
			UniquePatients: 0,
		},
	}
	cfg, err := extractionConfig.GetExtractionConfig(constants.ExtractionConfigFileName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load extraction Config",
		})
		return
	}

	for lineNum := 1; scanner.Scan(); lineNum++ {
		line := scanner.Bytes()

		var rawString string
		// Decode outer JSON string wrapper first
		if err := json.Unmarshal(line, &rawString); err != nil {
			log.Printf("Outer unmarshal failed (line %d): %v", lineNum, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid JSON string on line %d: %v", lineNum, err),
				"raw":   string(line),
			})
			return
		}

		// Split by newline – each should be an independent JSON object
		objects := strings.Split(rawString, "\n")
		for objNum, objLine := range objects {
			objLine = strings.TrimSpace(objLine)
			if objLine == "" {
				continue
			}

			dataQualityWarnings, err := apiFunctions.ProcessImportedJson(objLine, cfg)

			if err != nil {
				importReturn.ValidationErrors = append(importReturn.ValidationErrors, fmt.Sprintf("Error Processing on line %d object %d: %v", lineNum, objNum+1, err))
				continue
			}

			importReturn.TotalLinesProcessed++

			if len(dataQualityWarnings) > 0 {
				importReturn.DataQualityWarnings[""] = dataQualityWarnings
			} else {
				importReturn.RecordsImportedSuccessfully++
			}
		}
	}

	analytics, err := apiFunctions.ApiGetAnalyticsObject()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Encountered an issue, please try again later",
		})
		return
	}

	importReturn.Statistics.RecordsByType = analytics.RecordsByResourceType
	importReturn.Statistics.UniquePatients = analytics.NumberOfUniqueSubjects

	if err := scanner.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Error reading input: %v", err),
		})
		return
	}
	jsonBytes, err := json.Marshal(importReturn)

	c.JSON(http.StatusOK, gin.H{"result": string(jsonBytes)})
}

func apiGetRecordsById(c *gin.Context) {
	id := c.Param("id")

	fields := c.Query("fields")

	var fieldsArr []string

	if fields == "" {
		fieldsArr = nil
	} else {
		fieldsArr = strings.Split(fields, ",")
	}

	metadata, err := apiFunctions.GetRecordsById(id, fieldsArr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Encountered an issue, please try again later",
		})
		return
	}

	if metadata == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No Entry Found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": metadata})
}

func apiGetRecords(c *gin.Context) {
	fields := c.Query("fields")

	var fieldsArr []string

	if fields == "" {
		fieldsArr = nil
	} else {
		fieldsArr = strings.Split(fields, ",")
	}

	resourceType := c.Query("resourceType")

	subject := c.Query("subject")

	metadata, err := apiFunctions.GetRecords(resourceType, subject, fieldsArr)

	if err != nil {
		fmt.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Encountered an issue, please try again later",
		})
		return
	}

	if metadata == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No Entry Found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"metadata": metadata})
}

func apiPostTransform(c *gin.Context) {
	defer c.Request.Body.Close()

	jsonBytes, err := io.ReadAll(c.Request.Body)
	jsonString := string(jsonBytes)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	metadata, err := apiFunctions.TransformRequest(jsonString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if metadata == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No Entry Found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"metadata": metadata})
}

func apiGetAnalytics(c *gin.Context) {
	metadata, err := apiFunctions.ApiGetAnalytics()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Encountered an issue, please try again later",
		})
		return
	}

	if metadata == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No Entry Found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": metadata})
}
