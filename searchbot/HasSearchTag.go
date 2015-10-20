package searchbot

func HasSearchTag(tags []string) bool {
    for _, tag := range tags {
        if tag == "search" {
            return true
        }
    }
    return false
}

