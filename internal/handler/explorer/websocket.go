package explorer

import (
	"aurascan-backend/internal/websocket/pubsub"
	"ch-common-package/ginx"
	"ch-common-package/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func LatestHeightSub(c *gin.Context) {
	ws, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Infof("LatestHeightSub | failed to upgrade websocket protocol | err:%v", err)
		ginx.ResFailed(c, "Create Upgrade for websocket failed!")
		return
	}
	defer ws.Close()

	sub := pubsub.Subscribe()
	_ = ws.WriteJSON(map[string]string{"uid": sub.UUID().String()})
	for msg := range sub.C {
		err = ws.WriteJSON(msg)
		if err != nil {
			logger.Warnf("LatestHeightSub | failed to write message | err:%v", err)
			pubsub.Unsub(sub)
			break
		}
	}
}
