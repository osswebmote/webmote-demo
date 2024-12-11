package main

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/olahol/melody"
	"net/http"
	"webmote/handle"
)

func main() {
	e := echo.New()
	m := melody.New()
	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	e.Use(middleware.Static("./dist"))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ws/new", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"id": handle.NewId(),
		})
	})

	e.GET("/ws/:id", func(c echo.Context) error {
		id := c.Param("id")
		m.HandleRequestWithKeys(c.Response().Writer, c.Request(), map[string]any{"id": id})
		return nil
	})

	e.GET("/ws/game/:id", func(c echo.Context) error {
		id := c.Param("id")
		m.HandleRequestWithKeys(c.Response().Writer, c.Request(), map[string]any{"gid": id})
		return nil
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		ev := new(handle.Event)
		err := json.Unmarshal(msg, ev)
		if err != nil {
			e.Logger.Errorf("failed to unmarshal event: %v", err)
		}

		c, ret := handle.WS((s.Keys["id"]).(string), *ev)
		if !ret {
			return
		}

		r, err := json.Marshal(map[string]any{
			"e": ev.Event,
			"x": c.X,
			"y": c.Y,
		})
		if err != nil {
			e.Logger.Errorf("failed to marshal response: %v", err)
		}

		m.BroadcastFilter(r, func(session *melody.Session) bool {
			return s.Keys["id"] == session.Keys["gid"]
		})
	})

	m.HandleDisconnect(func(s *melody.Session) {
		if id, ok := s.Get("gid"); ok {
			sessions, err := m.Sessions()
			if err != nil {
				e.Logger.Errorf("failed to get sessions: %v", err)
				return
			}
			for _, session := range sessions {
				if session.Keys["id"] == id {
					session.Close()
				}
			}
			handle.Remove(id.(string))
		}
	})

	e.Logger.Fatal(e.StartTLS(":8000", "tls.cert", "tls.key"))
}
