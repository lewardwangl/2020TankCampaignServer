package ranks

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func WebSocketRanks(c *gin.Context) {
	conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
	if err != nil {

	}
	go func() {
		defer conn.Close()

		for {
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
			}
			fmt.Println(string(msg))
			err = wsutil.WriteServerMessage(conn, op, msg)
			if err != nil {
				// handle error
			}
		}
	}()
}
