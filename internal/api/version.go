package api

import (
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

// VersionHandler returns the git tag version.
func VersionHandler(c *gin.Context) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	out, err := cmd.Output()
	version := "dev"
	if err == nil {
		version = strings.TrimSpace(string(out))
	}
	c.JSON(http.StatusOK, gin.H{"version": version})
}
