package main

import (
	. "gopkg.in/check.v1"
)

func (s *MySuite) TestParsingSimpleStep(c *C) {
	parser := new(specParser)
	specText := SpecBuilder().specHeading("Spec heading with hash ").scenarioHeading("Scenario Heading").step("sample step").String()

	tokens, err := parser.generateTokens(specText)

	c.Assert(err, IsNil)
	c.Assert(len(tokens), Equals, 3)

	stepToken := tokens[2]
	c.Assert(stepToken.kind, Equals, stepKind)
	c.Assert(stepToken.value, Equals, "sample step")
}

func (s *MySuite) TestParsingEmptyStepTextShouldThrowError(c *C) {
	parser := new(specParser)
	specText := SpecBuilder().specHeading("Spec heading with hash ").scenarioHeading("Scenario Heading").step("").String()

	_, err := parser.generateTokens(specText)

	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "line no: 3, Step should not be blank")
}

func (s *MySuite) TestParsingStepWithParams(c *C) {
	parser := new(specParser)
	specText := SpecBuilder().specHeading("Spec heading with hash ").scenarioHeading("Scenario Heading").step("enter user \"john\"").String()

	tokens, err := parser.generateTokens(specText)

	c.Assert(err, IsNil)
	c.Assert(len(tokens), Equals, 3)

	stepToken := tokens[2]
	c.Assert(stepToken.kind, Equals, stepKind)
	c.Assert(stepToken.value, Equals, "enter user {static}")
	c.Assert(len(stepToken.args), Equals, 1)
	c.Assert(stepToken.args[0], Equals, "john")
}

func (s *MySuite) TestParsingStepWithParametersWithQuotes(c *C) {
	parser := new(specParser)
	specText := SpecBuilder().specHeading("Spec heading with hash ").scenarioHeading("Scenario Heading").step("\"param \\\"in quote\\\"\" step ").step("another * step with \"john 12 *-_{} \\\\ './;[]\" and \"second\"").String()

	tokens, err := parser.generateTokens(specText)

	c.Assert(err, IsNil)
	c.Assert(len(tokens), Equals, 4)

	firstStepToken := tokens[2]
	c.Assert(firstStepToken.kind, Equals, stepKind)
	c.Assert(firstStepToken.value, Equals, "{static} step")
	c.Assert(len(firstStepToken.args), Equals, 1)
	c.Assert(firstStepToken.args[0], Equals, "param \"in quote\"")

	secondStepToken := tokens[3]
	c.Assert(secondStepToken.kind, Equals, stepKind)
	c.Assert(secondStepToken.value, Equals, "another * step with {static} and {static}")
	c.Assert(len(secondStepToken.args), Equals, 2)
	c.Assert(secondStepToken.args[0], Equals, "john 12 *-_{} \\ './;[]")
	c.Assert(secondStepToken.args[1], Equals, "second")

}

func (s *MySuite) TestParsingStepWithUnmatchedOpeningQuote(c *C) {
	parser := new(specParser)
	specText := SpecBuilder().specHeading("Spec heading with hash ").scenarioHeading("Scenario Heading").step("sample step \"param").String()

	_, err := parser.generateTokens(specText)

	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "line no: 3, String not terminated")
}

func (s *MySuite) TestParsingStepWithEscaping(c *C) {
	parser := new(specParser)
	specText := SpecBuilder().specHeading("Spec heading with hash ").scenarioHeading("Scenario Heading").step("step with \\").String()

	tokens, err := parser.generateTokens(specText)

	c.Assert(err, IsNil)
	stepToken := tokens[2]
	c.Assert(stepToken.value, Equals, "step with")
}

func (s *MySuite) TestParsingExceptionIfStepContainsReservedChars(c *C) {
	parser := new(specParser)
	specText := SpecBuilder().specHeading("Spec heading with hash ").scenarioHeading("Scenario Heading").step("step with {braces}").String()

	_, err := parser.generateTokens(specText)

	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "line no: 3, '{' is a reserved character and should be escaped")
}

func (s *MySuite) TestParsingStepContainsEscapedReservedChars(c *C) {
	parser := new(specParser)
	specText := SpecBuilder().specHeading("Spec heading with hash ").scenarioHeading("Scenario Heading").step("step with \\{braces\\}").String()

	tokens, err := parser.generateTokens(specText)

	c.Assert(err, IsNil)
	stepToken := tokens[2]
	c.Assert(stepToken.value, Equals, "step with {braces}")
}

func (s *MySuite) TestParsingSimpleStepWithDynamicParameter(c *C) {
	parser := new(specParser)
	specText := SpecBuilder().specHeading("Spec heading with hash ").scenarioHeading("Scenario Heading").step("Step with \"static param\" and <name1>").String()

	tokens, err := parser.generateTokens(specText)
	c.Assert(err, IsNil)
	c.Assert(len(tokens), Equals, 3)

	stepToken := tokens[2]
	c.Assert(stepToken.value, Equals, "Step with {static} and {dynamic}")
	c.Assert(stepToken.args[0], Equals, "static param")
	c.Assert(stepToken.args[1], Equals, "name1")
}

func (s *MySuite) TestParsingStepWithUnmatchedDynamicParameterCharacter(c *C) {
	parser := new(specParser)
	specText := SpecBuilder().specHeading("Spec heading with hash ").scenarioHeading("Scenario Heading").step("Step with \"static param\" and <name1").String()

	_, err := parser.generateTokens(specText)

	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "line no: 3, Dynamic parameter not terminated")

}

func (s *MySuite) TestParsingContext(c *C) {
	parser := new(specParser)
	specText := SpecBuilder().specHeading("Spec heading with hash ").step("Context with \"param\"").scenarioHeading("Scenario Heading").String()

	tokens, err := parser.generateTokens(specText)

	c.Assert(err, IsNil)
	contextToken := tokens[1]
	c.Assert(contextToken.kind, Equals, stepKind)
	c.Assert(contextToken.value, Equals, "Context with {static}")
	c.Assert(contextToken.args[0], Equals, "param")
}

func (s *MySuite) TestParsingThrowsErrorWhenStepIsPresentWithoutStep(c *C) {
	parser := new(specParser)
	specText := SpecBuilder().step("step without spec heading").String()

	tokens, err := parser.generateTokens(specText)

	c.Assert(err, IsNil)
	c.Assert(tokens[0].kind, Equals, stepKind)
	c.Assert(tokens[0].value, Equals, "step without spec heading")

}

func (s *MySuite) TestParsingStepWithSimpleSpecialParameter(c *C) {
	parser := new(specParser)
	specText := SpecBuilder().specHeading("Spec heading with hash ").scenarioHeading("Scenario Heading").step("Step with special parameter <table:user.csv>").String()

	tokens, err := parser.generateTokens(specText)

	c.Assert(err, IsNil)
	c.Assert(len(tokens), Equals, 3)

	c.Assert(tokens[2].kind, Equals, stepKind)
	c.Assert(tokens[2].value, Equals, "Step with special parameter {special}")
	c.Assert(len(tokens[2].args), Equals, 1)
	c.Assert(tokens[2].args[0], Equals, "table:user.csv")
}

func (s *MySuite) TestParsingStepWithSpecialParametersWithWhiteSpaces(c *C) {
	parser := new(specParser)
	specText := SpecBuilder().specHeading("Spec heading with hash ").step("Step with \"first\" and special parameter <table : user.csv>").step("Another with <name> and <file  :something.txt>").String()

	tokens, err := parser.generateTokens(specText)

	c.Assert(err, IsNil)
	c.Assert(len(tokens), Equals, 3)

	c.Assert(tokens[1].kind, Equals, stepKind)
	c.Assert(tokens[1].value, Equals, "Step with {static} and special parameter {special}")
	c.Assert(len(tokens[1].args), Equals, 2)
	c.Assert(tokens[1].args[0], Equals, "first")
	c.Assert(tokens[1].args[1], Equals, "table : user.csv")

	c.Assert(tokens[2].kind, Equals, stepKind)
	c.Assert(tokens[2].value, Equals, "Another with {dynamic} and {special}")
	c.Assert(len(tokens[2].args), Equals, 2)
	c.Assert(tokens[2].args[0], Equals, "name")
	c.Assert(tokens[2].args[1], Equals, "file  :something.txt")
}
