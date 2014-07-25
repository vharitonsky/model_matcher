package lib

import (
	"testing"
)

func TestCleanName(t *testing.T) {
	name := "!\"#$%&'()*+,-./:;<=>?@[]^_`{|}~"
	clean_name := CleanName(name)
	if clean_name != " " {
		t.Error("Expected clean name to remove all characters,", clean_name, "left")
	}
}

func TestSplitName(t *testing.T) {
	name := "three words NaMe"
	name_map := SplitName(name)
	_, found := name_map["three"]
	if !found {
		t.Error("Expected split name result to have 'three' key")
	}
	_, found = name_map["words"]
	if !found {
		t.Error("Expected split name result to have 'words' key")
	}
	_, found = name_map["NaMe"]
	if !found {
		t.Error("Expected split name result to have 'NaMe' key")
	}
}

func TestMatchNames(t *testing.T) {
	name_map_a := map[string]bool{"this": true, "is": true, "a": true, "long": true, "name": true}
	name_map_b := map[string]bool{"this": true, "is": true, "name": true}
	match := MatchNames(name_map_a, name_map_b)
	if !match {
		t.Error("Expected names to match", name_map_a, name_map_b)
	}
	name_map_b = map[string]bool{"this": true, "is": true, "not": true, "name": true}
	match = MatchNames(name_map_a, name_map_b)
	if match {
		t.Error("Expected names to not match", name_map_a, name_map_b)
	}
}
