package util

import "strings"

func ScopeShouldContain(scopes []string, scope string) bool {
	for _, s := range scopes {
		if res := strings.Contains(scope, s); !res {
			return false
		}
	}
	return true
}

func ScopeAddTag(currentScope, newTag string) string {
	if currentScope == "" {
		return newTag
	}
	var newScope = currentScope
	if !ScopeShouldContain([]string{newScope}, newTag) {
		newScope = newScope + " " + newTag
	}
	return newScope
}
