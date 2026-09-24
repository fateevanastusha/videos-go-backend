package validators

func ValidateTitleAuthor(title, author string) []ErrorMessage {
	errs := []ErrorMessage{}
	if len(title) > 40 || len(title) == 0 {
		errs = append(errs, ErrorMessage{Message: "must be no longer 40 and not be empty", Field: "title"})
	}
	if len(author) > 20 || len(author) == 0 {
		errs = append(errs, ErrorMessage{Message: "must be no longer 20 and not be empty", Field: "author"})
	}
	return errs
}
