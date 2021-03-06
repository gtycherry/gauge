package main

import (
	"fmt"
	. "gopkg.in/check.v1"
)

type specBuilder struct {
	lines []string
}

func SpecBuilder() *specBuilder {
	return &specBuilder{lines: make([]string, 0)}
}

func (specBuilder *specBuilder) addPrefix(prefix string, line string) string {
	return fmt.Sprintf("%s%s\n", prefix, line)
}

func (specBuilder *specBuilder) String() string {
	var result string
	for _, line := range specBuilder.lines {
		result = fmt.Sprintf("%s%s", result, line)
	}
	return result
}

func (specBuilder *specBuilder) specHeading(heading string) *specBuilder {
	line := specBuilder.addPrefix("#", heading)
	specBuilder.lines = append(specBuilder.lines, line)
	return specBuilder
}

func (specBuilder *specBuilder) scenarioHeading(heading string) *specBuilder {
	line := specBuilder.addPrefix("##", heading)
	specBuilder.lines = append(specBuilder.lines, line)
	return specBuilder
}

func (specBuilder *specBuilder) step(stepText string) *specBuilder {
	line := specBuilder.addPrefix("* ", stepText)
	specBuilder.lines = append(specBuilder.lines, line)
	return specBuilder
}

func (specBuilder *specBuilder) tags(tags ...string) *specBuilder {
	tagText := ""
	for i, tag := range tags {
		tagText = fmt.Sprintf("%s%s", tagText, tag)
		if i != len(tags)-1 {
			tagText = fmt.Sprintf("%s,", tagText)
		}
	}
	line := specBuilder.addPrefix("tags: ", tagText)
	specBuilder.lines = append(specBuilder.lines, line)
	return specBuilder
}

func (specBuilder *specBuilder) tableHeader(cells ...string) *specBuilder {
	return specBuilder.tableRow(cells...)
}
func (specBuilder *specBuilder) tableRow(cells ...string) *specBuilder {
	rowInMarkdown := "|"
	for _, cell := range cells {
		rowInMarkdown = fmt.Sprintf("%s%s|", rowInMarkdown, cell)
	}
	specBuilder.lines = append(specBuilder.lines, fmt.Sprintf("%s\n", rowInMarkdown))
	return specBuilder
}

func (specBuilder *specBuilder) text(comment string) *specBuilder {
	specBuilder.lines = append(specBuilder.lines, fmt.Sprintf("%s\n", comment))
	return specBuilder
}

func (s *MySuite) TestIsInState(c *C) {
	c.Assert(isInState(1, 1), Equals, true)
	c.Assert(isInState(1, 3, 2), Equals, true)
	c.Assert(isInState(4, 1, 2), Equals, false)
	c.Assert(isInState(4), Equals, false)
}

func (s *MySuite) TestIsInAnyState(c *C) {
	c.Assert(isInAnyState(4, 4), Equals, true)
	c.Assert(isInAnyState(4, 1, 4), Equals, true)
	c.Assert(isInAnyState(8, 1, 3), Equals, false)
	c.Assert(isInAnyState(8), Equals, false)
}

func (s *MySuite) TestRetainStates(c *C) {
	oldState := 5
	retainStates(&oldState, 1)
	c.Assert(oldState, Equals, 1)

	oldState = 5
	retainStates(&oldState, 4, 1)
	c.Assert(oldState, Equals, 5)

	oldState = 8
	retainStates(&oldState, 4, 6)
	c.Assert(oldState, Equals, 0)
}

func (s *MySuite) TestAddStates(c *C) {
	oldState := 4
	addStates(&oldState, 1)
	c.Assert(oldState, Equals, 5)

	oldState = 8
	addStates(&oldState, 1, 2)
	c.Assert(oldState, Equals, 11)

	oldState = 8
	addStates(&oldState)
	c.Assert(oldState, Equals, 8)
}

func (s *MySuite) TestAreUnderlinedForEmptyArray(c *C) {
	emptyAray := make([]string, 0)
	c.Assert(false, Equals, areUnderlined(emptyAray))
}

func (s *MySuite) TestSpecBuilderSpecHeading(c *C) {
	heading := SpecBuilder().specHeading("spec heading").String()
	c.Assert(heading, Equals, "#spec heading\n")
}

func (s *MySuite) TestSpecBuilderScenarioHeading(c *C) {
	heading := SpecBuilder().scenarioHeading("scenario heading").String()
	c.Assert(heading, Equals, "##scenario heading\n")
}

func (s *MySuite) TestSpecBuilderStep(c *C) {
	step := SpecBuilder().step("sample step").String()
	c.Assert(step, Equals, "* sample step\n")
}

func (s *MySuite) TestSpecBuilderTags(c *C) {
	tags := SpecBuilder().tags("tag1", "tag2").String()
	c.Assert(tags, Equals, "tags: tag1,tag2\n")
}

func (s *MySuite) TestSpecBuilderWithFreeText(c *C) {
	freeText := SpecBuilder().text("some free text").String()
	c.Assert(freeText, Equals, "some free text\n")
}

func (s *MySuite) TestSpecBuilderWithSampleSpec(c *C) {
	spec := SpecBuilder().specHeading("spec heading").tags("tag1", "tag2").step("context here").scenarioHeading("scenario heading").text("comment").step("sample step").scenarioHeading("scenario 2").step("step 2")
	c.Assert(spec.String(), Equals, "#spec heading\ntags: tag1,tag2\n* context here\n##scenario heading\ncomment\n* sample step\n##scenario 2\n* step 2\n")
}
