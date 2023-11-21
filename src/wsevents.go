package src

import (
	"bytes"
	"html/template"
)

// register client
// unregister client

// change client username
func NewUserNameHTMl(p Player) (string, error) {
	tmpl, err := template.ParseFiles("./view/player.html")
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, p); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

// wait for other clients to connect
// start game
// wait turn
// play turn
