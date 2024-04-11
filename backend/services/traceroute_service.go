package services

import (
	"context"
	"devtools/backend/traceroute"
	"devtools/backend/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"sync"
)

type TraceRoute struct {
	Routes   []*traceroute.TraceResult
	Msg      string
	Success  bool
	Complete bool
}

func (t *TraceRoute) addRoute(route *traceroute.TraceResult) {
	t.Routes = append(t.Routes, route)
}

type traceRouteService struct {
	ctx        context.Context
	traceCache map[string]*TraceRoute
}

var traceRouter *traceRouteService
var onceTraceRouter sync.Once

func TraceRouter() *traceRouteService {
	if traceRouter == nil {
		onceTraceRouter.Do(func() {
			traceRouter = &traceRouteService{
				traceCache: make(map[string]*TraceRoute),
			}
		})
	}
	return traceRouter
}

func (t *traceRouteService) Start(ctx context.Context) {
	t.ctx = ctx
}

func (t *traceRouteService) GetTraceRoute(dest string) (resp types.JSResp) {
	var traceResult *TraceRoute
	traceResult, ok := t.traceCache[dest]
	if !ok {
		traceResult = &TraceRoute{
			Success: true,
		}
		t.traceCache[dest] = traceResult
		go func() {
			t.startTraceRoute(dest, traceResult)
		}()
	}
	if !traceResult.Success {
		resp.Msg = traceResult.Msg
	}
	resp.Data = map[string]any{
		"routes":   traceResult.Routes,
		"complete": traceResult.Complete,
	}
	resp.Success = true
	return
}

func (t *traceRouteService) startTraceRoute(dest string, routeResult *TraceRoute) {
	routes, errsChan, err := traceroute.Traceroute(dest)
	if err != nil {
		routeResult.Success = false
		routeResult.Msg = err.Error()
		return
	}
	runtime.LogDebugf(t.ctx, "start trace route dest: %v", dest)

	for {
		select {
		case err := <-errsChan:
			runtime.LogDebugf(t.ctx, "trace route err: %v", err)
			routeResult.Msg = err.Error()
			return
		case route, ok := <-routes:
			if !ok {
				routeResult.Complete = true
				return
			}
			runtime.LogDebugf(t.ctx, "get trace route : %v", route)
			routeResult.addRoute(route)
		case <-t.ctx.Done():
			routeResult.Complete = true
			return
		}
	}
}
