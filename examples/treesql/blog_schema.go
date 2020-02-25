package treesql

var BlogSchema = &SchemaDesc{
	Tables: map[string]*TableDesc{
		"posts": {
			Name: "posts",
			Columns: map[string]*ColDesc{
				"id":    {Name: "id"},
				"body":  {Name: "body"},
				"title": {Name: "title"},
				// vv here so that both have a col that starts with p
				"pics": {Name: "pics"},
			},
		},
		"comments": {
			Name: "comments",
			Columns: map[string]*ColDesc{
				"id":      {Name: "id"},
				"body":    {Name: "body"},
				"post_id": {Name: "post_id"},
			},
		},
	},
}
