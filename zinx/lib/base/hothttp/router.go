/**************************************************************************************
Code Description    : 路由
Code Vesion         :
					|------------------------------------------------------------|
						  Version    					Editor            Time
							1.0        					yuansudong        2016.4.12
					|------------------------------------------------------------|
Version Description	:
                    |------------------------------------------------------------|
						  Version
							1.0
								 ....
					|------------------------------------------------------------|
***************************************************************************************/

package hothttp

import (
	"fmt"

	"bangseller.com/lib/base/util"
)

type (
	router struct {
		store []map[string]*handle
	}
	handle struct {
		method        string
		executeHandle InterfaceHandle
	}
)

// newRouter 用于新建一个路由
func newRouter() *router {
	pRouter := &router{
		store: make([]map[string]*handle, routerBucket),
	}
	for i := 0; i < routerBucket; i++ {
		pRouter.store[i] = make(map[string]*handle)
	}
	return pRouter
}

// RegisterInterface 用于添加行为码
func (r *router) RegisterInterface(method string, uri string, function InterfaceHandle) {
	valueCode := util.StringToInt(uri)
	_, isExists := r.store[valueCode%routerBucket][uri]
	if isExists {
		panic(fmt.Sprintf("the interface %s already isExists", uri))
	}
	r.store[valueCode%routerBucket][uri] = &handle{
		method:        method,
		executeHandle: function,
	}
}

// printInfo 用于打印接口信息
func (r *router) printInfo() {

}
