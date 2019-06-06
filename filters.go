package event

type subscriptionFilter func(Event) bool

func FilterAll() subscriptionFilter {
	return func(e Event) bool {
		return true
	}
}

func FilterName(s []string) subscriptionFilter {
	return func(e Event) bool {
		for _, str := range s {
			if e.Name() == str {
				return true
			}
		}
		return false
	}
}

func FilterNotName(s []string) subscriptionFilter {
	return func(e Event) bool {
		switch FilterName(s)(e) {
		case true:
			return false
		default:
			return true
		}
	}
}

func FilterEveryOther() subscriptionFilter {
	i := 1
	return func(e Event) bool {
		i++
		if i%2 == 0 {
			return true
		}
		return false
	}
}
