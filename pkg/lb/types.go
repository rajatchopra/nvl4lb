
package lb

import (
	"net/http"
)


var (
	VirtInterface string
)

type server struct {
	http.Server
}
