package main

import (
	"strings"
)

/*
* Used for checking if a topic is valid for publishing/subscribing against a given allowed topic.
 */
func MatchTopicAgainst(topic, against string) bool {
	words1 := strings.Split(topic, "/")
	words2 := strings.Split(against, "/")
	return matchTokens(words1, words2)
}

func matchTokens(words1, words2 []string) bool {
	if len(words1) == 0 && len(words2) == 0 {
		return true
	}
	if len(words2) == 0 || len(words1) == 0 {
		return false
	}
	if words2[0] == "#" {
		return true
	}
	if words2[0] == "+" || words1[0] == words2[0] {
		if len(words1) == 1 {
			words1 = make([]string, 0)
		} else {
			words1 = words1[1:len(words1)]
		}
		if len(words2) == 1 {
			words2 = make([]string, 0)
		} else {
			words2 = words2[1:len(words2)]
		}
		return matchTokens(words1, words2)
	}
	return false
}
func IsTopicAllowed(topic string, topics []string) bool {
	for i := 0; i < len(topics); i++ {
		if MatchTopicAgainst(topic, topics[i]) {
			return true
		}
	}
	return false
}

/*
func test() {
	topic := "leval1/level2/level3/level4"
	topics := []string{"level1/level2", "level1/level2/level3/level4", "level1/level2/+/level4", "level1/#", "level2/#"}
	results := make([]string, 0)
	for i := 1; i < len(topics); i += 1 {
		if matched := MatchTopicAgainst(topic, topics[i]); matched {
			results = append(results, topics[i])
		}
	}

	for i := 1; i < len(results); i += 1 {
		fmt.Println(results[i])
	}

}
*/
