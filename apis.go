package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	dbObj, err := NewDatabase(SqlDBName)
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

	// Define a simple GET endpoint
	r.GET("/ping", getPong)
	r.POST("/import", postImport)
	r.GET("/GetByID/:id", GetById)

	return r
}

// api tester
func getPong(c *gin.Context) {
	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// adds items from JSONL received in the request body.
func postImport(c *gin.Context) {
	scanner := bufio.NewScanner(c.Request.Body)
	defer c.Request.Body.Close()

	importReturn := ImportReturn{
		totalLinesProcessed:         0,
		recordsImportedSuccessfully: 0,
		validationErrors:            make([]string, 0),
		dataQualityWarnings:         make(map[string][]string),
		statistics: Statistics{
			RecordsByType:  make(map[string]int),
			uniquePatients: 0,
		},
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

			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(objLine), &obj); err != nil {
				log.Printf("Inner unmarshal failed (line %d obj %d): %v", lineNum, objNum+1, err)
				importReturn.validationErrors = append(importReturn.validationErrors, fmt.Sprintf("Invalid inner JSON on line %d object %d: %v", lineNum, objNum+1, err))
				continue
			}

			dataQualityWarnings, err := ProcessImportedJson(obj)

			if err != nil {
				importReturn.validationErrors = append(importReturn.validationErrors, fmt.Sprintf("Error Processing on line %d object %d: %v", lineNum, objNum+1, err))
				continue
			}

			importReturn.totalLinesProcessed++

			if len(dataQualityWarnings) > 0 {
				importReturn.dataQualityWarnings[""] = dataQualityWarnings
			}

		}
	}

	if err := scanner.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Error reading input: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Import complete"})
}

func GetById(c *gin.Context) {
	id := c.Param("id")

	db, err := Connect(SqlDBName)

	if err != nil {

	}

	defer db.Close()

	metadata, err := GetResourceById(db, id)

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

	c.JSON(http.StatusOK, gin.H{"metadata": metadata})
}
