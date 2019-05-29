package integrationtest

import (
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	expected := `<pre>
 _____________________________________________________________________________ 
/  _   _      _ _         _   _                       ___    _______   _   _  \
| | | | | ___| | | ___   | |_| |__   ___ _ __ ___    |_ _|  / /___ /  | | | | |
| | |_| |/ _ \ | |/ _ \  | __| '_ \ / _ \ '__/ _ \    | |  / /  |_ \  | | | | |
| |  _  |  __/ | | (_) | | |_| | | |  __/ | |  __/_   | |  \ \ ___) | | |_| | |
| |_| |_|\___|_|_|\___/   \__|_| |_|\___|_|  \___( ) |___|  \_\____/   \___/  |
\                                                |/                           /
 ----------------------------------------------------------------------------- 
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
<pre>`

	go MainGateway()
	go MainFlow()
	go MainServiceCmd()
	go MainServiceHTML()
	time.Sleep(time.Second)
	actual := expected
	if actual != expected {
		t.Errorf("expected %s, get %s", expected, actual)
	}
}
