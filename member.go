package valourgo

import "net/http"

func (n *Node) MyMember(planetID PlanetID) (*Member, error) {
	me, err := n.Me()

	if err != nil {
		return nil, err
	}

	member, err := n.MemberByUser(planetID, me.ID)

	if err != nil {
		return nil, err
	}

	return member, nil
}

func (n *Node) Member(id MemberID) (*Member, error) {
	var member Member

	if err := n.requestJSON(http.MethodGet, "api/members/"+id.String(), nil, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

func (n *Node) MemberByUser(planetID PlanetID, id UserID) (*Member, error) {
	var member Member

	if err := n.requestJSON(http.MethodGet, "api/members/byuser/"+planetID.String()+"/"+id.String(), nil, &member); err != nil {
		return nil, err
	}

	return &member, nil
}
