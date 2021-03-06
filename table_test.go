package main

import (
	. "gopkg.in/check.v1"
)

func (s *MySuite) TestIsInitalized(c *C) {
	var table table
	c.Assert(table.isInitialized(), Equals, false)
	c.Assert(table.getRowCount(), Equals, 0)

	table.addHeaders([]string{"one", "two", "three"})

	c.Assert(table.isInitialized(), Equals, true)
}

func (s *MySuite) TestShouldAddHeaders(c *C) {
	var table table

	table.addHeaders([]string{"one", "two", "three"})

	c.Assert(len(table.headers), Equals, 3)
	c.Assert(table.headerIndexMap["one"], Equals, 0)
	c.Assert(table.headers[0], Equals, "one")
	c.Assert(table.headerIndexMap["two"], Equals, 1)
	c.Assert(table.headers[1], Equals, "two")
	c.Assert(table.headerIndexMap["three"], Equals, 2)
	c.Assert(table.headers[2], Equals, "three")
}

func (s *MySuite) TestShouldAddRows(c *C) {
	var table table

	table.addHeaders([]string{"one", "two", "three"})
	table.addRowValues([]string{"foo", "bar", "baz"})
	table.addRowValues([]string{"john", "jim"})

	c.Assert(table.getRowCount(), Equals, 2)
	column1 := table.get("one")
	c.Assert(len(column1), Equals, 2)
	c.Assert(column1[0].value, Equals, "foo")
	c.Assert(column1[0].cellType, Equals, static)
	c.Assert(column1[1].value, Equals, "john")
	c.Assert(column1[1].cellType, Equals, static)

	column2 := table.get("two")
	c.Assert(len(column2), Equals, 2)
	c.Assert(column2[0].value, Equals, "bar")
	c.Assert(column2[0].cellType, Equals, static)
	c.Assert(column2[1].value, Equals, "jim")
	c.Assert(column2[1].cellType, Equals, static)

	column3 := table.get("three")
	c.Assert(len(column3), Equals, 2)
	c.Assert(column3[0].value, Equals, "baz")
	c.Assert(column3[0].cellType, Equals, static)
	c.Assert(column3[1].value, Equals, "")
	c.Assert(column3[1].cellType, Equals, static)

}

func (s *MySuite) TestCoulmnNameExists(c *C) {
	var table table

	table.addHeaders([]string{"one", "two", "three"})
	table.addRowValues([]string{"foo", "bar", "baz"})
	table.addRowValues([]string{"john", "jim", "jack"})

	c.Assert(table.headerExists("one"), Equals, true)
	c.Assert(table.headerExists("two"), Equals, true)
	c.Assert(table.headerExists("four"), Equals, false)

}

func (s *MySuite) TestGetInvalidColumn(c *C) {
	var table table

	table.addHeaders([]string{"one", "two", "three"})
	table.addRowValues([]string{"foo", "bar", "baz"})
	table.addRowValues([]string{"john", "jim", "jack"})

	c.Assert(func() { table.get("four") }, Panics, "Table column four not found")
}

func (s *MySuite) TestGetRows(c *C) {
	var table table

	table.addHeaders([]string{"one", "two", "three"})
	table.addRowValues([]string{"foo", "bar", "baz"})
	table.addRowValues([]string{"john", "jim", "jack"})

	rows := table.getRows()
	c.Assert(len(rows), Equals, 2)
	firstRow := rows[0]
	c.Assert(firstRow[0], Equals, "foo")
	c.Assert(firstRow[1], Equals, "bar")
	c.Assert(firstRow[2], Equals, "baz")

	secondRow := rows[1]
	c.Assert(secondRow[0], Equals, "john")
	c.Assert(secondRow[1], Equals, "jim")
	c.Assert(secondRow[2], Equals, "jack")
}
