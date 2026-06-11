package main

import (
	"crypto/rand"
	"runtime"
	"sync"
	"unsafe"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// sessionLogRegistry assigns stable log IDs when the MCP SDK does not expose a
// transport session ID (SSE and stdio connections return an empty SessionID()).
type sessionLogRegistry struct {
	ids sync.Map // uintptr -> string
}

func (r *sessionLogRegistry) resolve(sess mcp.Session) string {
	if sess == nil {
		return ""
	}
	if id := sess.ID(); id != "" {
		return id
	}

	ss, ok := sess.(*mcp.ServerSession)
	if !ok {
		return rand.Text()
	}

	key := sessionPtr(ss)
	if v, loaded := r.ids.Load(key); loaded {
		return v.(string)
	}

	id := rand.Text()
	if v, loaded := r.ids.LoadOrStore(key, id); loaded {
		return v.(string)
	}

	runtime.SetFinalizer(ss, func(s *mcp.ServerSession) {
		r.ids.Delete(sessionPtr(s))
	})
	return id
}

func sessionPtr(ss *mcp.ServerSession) uintptr {
	return uintptr(unsafe.Pointer(ss))
}
