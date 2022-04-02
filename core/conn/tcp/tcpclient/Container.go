package tcpclient

import "games/core/conn"

/// <summary>
/// sessions 客户端容器
/// <summary>
var sessions = conn.NewSessions()

func Get(id int64) conn.Session {
	return sessions.Get(id)
}

func Count() int {
	return sessions.Count()
}

func CloseAll() {
	sessions.CloseAll()
}

func Wait() {
	sessions.Wait()
}

func Stop() {
	sessions.Stop()
}
