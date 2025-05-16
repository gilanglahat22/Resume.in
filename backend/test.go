// This file is used to ensure Go can resolve and import all dependencies
package main

import (
	// Standard library imports
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	// Third-party imports
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
)

// dummyFunc is to avoid go compiler complaining about unused imports
func dummyFunc() {
	_ = context.Background()
	_ = sql.ErrNoRows
	_ = json.Marshal
	_ = fmt.Sprintf
	_ = os.Getenv
	_ = time.Now()
	_ = gin.Default()
	_ = uuid.New()
	_ = &pq.Driver{}
	_ = pgvector.NewVector([]float32{0.1, 0.2})
} 