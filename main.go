package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Label Node
	r.POST("/label", func(c *gin.Context) {
		// Parse JSON
		var json struct {
			NodeID      string `json:"NodeID" binding:"required"`
			ServiceName string `json:"ServiceName" binding:"required"`
		}

		if c.Bind(&json) == nil {
			err := SystemExec([][]string{
				{
					"docker", "node", "update", "--label-add", fmt.Sprintf("service.name=%s", json.ServiceName), json.NodeID,
				},
			})

			if err != nil {
				panic(fmt.Errorf("Command failed with: %+v", err))

				c.JSON(http.StatusOK, gin.H{"message": "Failed to label node"})
			}

			c.JSON(http.StatusOK, gin.H{"message": "Node labeled successfully"})
		}
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":3300")
}

// SystemExec execute a command in the system
func SystemExec(commands [][]string) error {
	for _, c := range commands {
		cmd := exec.Command(c[0], c[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
