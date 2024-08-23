package explorer

import (
	"aurascan-backend/internal/config"
	"aurascan-backend/model/schema"
	"ch-common-package/ginx"
	"ch-common-package/logger"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"os"
	"os/exec"
	"time"
)

func CompileLeo(c *gin.Context) {
	var lsc schema.LeoSourceCode
	if err := c.BindJSON(&lsc); err != nil {
		logger.Errorf("CompileLeo BindJSON | %v", err)
		ginx.ResFailed(c, "missing parameter")
		return
	}

	s, err := LeoBuild(lsc)
	if err != nil {
		logger.Error("build leoBuild | %v", err)
		ginx.ResFailed(c, "build failed")
		return
	}
	ginx.ResSuccess(c, s)
}

func LeoBuild(lsc schema.LeoSourceCode) (string, error) {
	reg, err := regexp2.Compile(`(?<=program )\w+(?=\.aleo {)`, 0)
	if err != nil {
		return "", fmt.Errorf("compile | %v", err)
	}
	m, err := reg.FindStringMatch(lsc.Raw)
	if m == nil {
		return "", fmt.Errorf("findStringMatch | get null string")
	}
	if err != nil {
		logger.Errorf("findStringMatch | %v", err)
		return "", err
	}

	pargramName := m.String()
	pargramHome := fmt.Sprintf("%s/%s-%d", config.Global.LeoCompile.HomePath, pargramName, time.Now().UnixNano())
	err = os.Mkdir(pargramHome, 0755)
	if err != nil {
		return "", fmt.Errorf("mkdir | failed to create directory for %s | %v", pargramName, err)
	}

	cmd := exec.Command(config.Global.LeoCompile.BinPath, "new", pargramName)
	cmd.Dir = pargramHome
	_, err = cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("combinedOutput | failed to new program | %v", err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/%s/src/main.leo", pargramHome, pargramName), []byte(lsc.Raw), 0644)
	if err != nil {
		return "", fmt.Errorf("writeFile | failed to replace main.leo | %v", err)
	}

	cmd = exec.Command(config.Global.LeoCompile.BinPath, "build")
	cmd.Dir = fmt.Sprintf("%s/%s", pargramHome, pargramName)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("combinedOutput | failed to build program | %v", err)
	}

	b, err := os.ReadFile(fmt.Sprintf("%s/build/main.aleo", cmd.Dir))
	if err != nil {
		return "", fmt.Errorf("readFile | failed to ReadFile | %v", err)
	}
	return string(b), nil
}
