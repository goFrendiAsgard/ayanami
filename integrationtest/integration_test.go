package integrationtest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	expected := fmt.Sprintf("%s%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s",
		"<pre>",
		` _____________________________________________________________________________ `,
		`/  _   _      _ _         _   _                       ___    _______   _   _  \`,
		`| | | | | ___| | | ___   | |_| |__   ___ _ __ ___    |_ _|  / /___ /  | | | | |`,
		`| | |_| |/ _ \ | |/ _ \  | __| '_ \ / _ \ '__/ _ \    | |  / /  |_ \  | | | | |`,
		`| |  _  |  __/ | | (_) | | |_| | | |  __/ | |  __/_   | |  \ \ ___) | | |_| | |`,
		`| |_| |_|\___|_|_|\___/   \__|_| |_|\___|_|  \___( ) |___|  \_\____/   \___/  |`,
		`|                                                |/                           |`,
		`\                                                                             /`,
		` ----------------------------------------------------------------------------- `,
		`        \   ^__^`,
		`         \  (oo)\_______`,
		`            (__)\       )\/\`,
		`                ||----w |`,
		`                ||     ||`,
		"</pre>",
	)
	// run services
	go MainGateway()
	go MainFlow()
	go MainServiceCmd()
	go MainServiceHTML()
	// wait for two seconds
	time.Sleep(100 * time.Millisecond)
	// emulate request
	response, err := http.Get(fmt.Sprintf("http://localhost:8080/?text=%s", url.QueryEscape("Hello there, I <3 U")))
	if err != nil {
		t.Errorf("Get error %s", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Get error %s", err)
	}
	actual := string(body)
	if actual != expected {
		t.Errorf("expected :\n%s, get :\n%s", expected, actual)
	}
}
