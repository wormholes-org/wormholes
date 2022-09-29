package sgiccommon

import "os"

var Sigc = make(chan os.Signal, 1)
