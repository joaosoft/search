package dbr

import (
	"fmt"
)

type functionAgg struct {
	function    string
	arguments   []interface{}
	orders      orders
	filters     *conditions
	withinGroup bool

	*functionBase
}

func newFunctionAgg(function string, arguments ...interface{}) *functionAgg {
	return &functionAgg{functionBase: newFunctionBase(false, false), function: function, arguments: arguments, filters: newConditions(nil)}
}

func (c *functionAgg) OrderBy(column string, direction direction) *functionAgg {
	c.orders = append(c.orders, &order{column: column, direction: direction})
	return c
}

func (c *functionAgg) OrderAsc(columns ...string) *functionAgg {
	for _, column := range columns {
		c.orders = append(c.orders, &order{column: column, direction: OrderAsc})
	}
	return c
}

func (c *functionAgg) WithinGroup(enable bool) *functionAgg {
	c.withinGroup = enable
	return c
}

func (c *functionAgg) OrderDesc(columns ...string) *functionAgg {
	for _, column := range columns {
		c.orders = append(c.orders, &order{column: column, direction: OrderDesc})
	}
	return c
}

func (c *functionAgg) Filter(query interface{}, values ...interface{}) *functionAgg {
	c.filters.list = append(c.filters.list, &condition{operator: OperatorAnd, query: query, values: values})
	return c
}

func (c *functionAgg) Build(db *db) (string, error) {
	c.db = db

	functionBase := newFunctionBase(false, false, db)

	var arguments string

	lenArgs := len(c.arguments)
	for i, argument := range c.arguments {
		expression, err := handleBuild(functionBase, argument)
		if err != nil {
			return "", err
		}

		arguments += expression

		if i < lenArgs-1 {
			arguments += ", "
		}
	}

	orders, err := c.orders.Build()
	if err != nil {
		return "", err
	}

	filters, err := c.filters.Build(c.db)
	if err != nil {
		return "", err
	}

	if len(filters) > 0 {
		filters = fmt.Sprintf(" %s (%s %s)", constFunctionFilter, constFunctionWhere, filters)
	}

	if c.withinGroup {
		return fmt.Sprintf("%s(%s) %s (%s)%s", c.function, arguments, constFunctionWithinGroup, orders, filters), nil
	}

	return fmt.Sprintf("%s(%s%s)%s", c.function, arguments, orders, filters), nil
}
