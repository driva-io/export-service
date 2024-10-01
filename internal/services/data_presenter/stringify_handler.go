package data_presenter

func handleStringify(source map[string]any, location any) (any, error) {
	stringify := map[string]any{
		"$prop": location,
	}

	return handleJoinBy(source, stringify)
}
