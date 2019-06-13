package treesql

var BlogSchema = &SchemaDesc{
	Tables: map[string]*TableDesc{
		"posts": {
			Columns: map[string]*ColDesc{
				"id":    {},
				"body":  {},
				"title": {},
				// vv here so that both have a col that starts with p
				"pics": {},
			},
		},
		"comments": {
			Columns: map[string]*ColDesc{
				"id":      {},
				"body":    {},
				"post_id": {},
			},
		},
	},
}
