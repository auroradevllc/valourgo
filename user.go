package valourgo

import "net/http"

func (n *Node) Me() (*User, error) {
	if n.me != nil {
		return n.me, nil
	}

	var user User

	if err := n.requestJSON(http.MethodGet, "api/users/me", nil, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
