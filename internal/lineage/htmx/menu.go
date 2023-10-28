package htmx

type MenuItem struct {
	Key  string
	Text string
	Href string
	Role string
	Icon string
}

var MenuItems = []MenuItem{
	{Key: "datasets", Text: "Datasets", Href: "/lineage/datasets", Icon: "table"},
	{Key: "events", Text: "Events", Href: "/lineage/requests", Icon: "list-alt"},
	{Key: "jobs", Text: "Jobs", Href: "/lineage/jobs", Icon: "cogs"},
}

func BuildMenuItems(chosenKey string) []MenuItem {
	var res []MenuItem
	for _, it := range MenuItems {
		if it.Key == chosenKey {
			it.Role = "button"
		}
		res = append(res, it)
	}
	return res
}
