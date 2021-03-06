// Copyright 2016 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package xeval

import (
	. "github.com/pingcap/check"
	"github.com/pingcap/tidb/util/types"
	"github.com/pingcap/tipb/go-tipb"
)

func (s *testEvalSuite) TestEvalCaseWhen(c *C) {
	colID := int64(1)
	row := make(map[int64]types.Datum)
	row[colID] = types.NewIntDatum(100)
	xevaluator := &Evaluator{Row: row}
	trueCond := types.NewIntDatum(1)
	falseCond := types.NewIntDatum(0)
	nullCond := types.Datum{}
	nullCond.SetNull()
	cases := []struct {
		expr   *tipb.Expr
		result types.Datum
	}{
		{
			expr: buildExpr(tipb.ExprType_Case,
				falseCond, types.NewStringDatum("case1"),
				trueCond, types.NewStringDatum("case2"),
				trueCond, types.NewStringDatum("case3")),
			result: types.NewStringDatum("case2"),
		},
		{
			expr: buildExpr(tipb.ExprType_Case,
				falseCond, types.NewStringDatum("case1"),
				falseCond, types.NewStringDatum("case2"),
				falseCond, types.NewStringDatum("case3"),
				types.NewStringDatum("Else")),
			result: types.NewStringDatum("Else"),
		},
		{
			expr: buildExpr(tipb.ExprType_Case,
				falseCond, types.NewStringDatum("case1"),
				falseCond, types.NewStringDatum("case2"),
				falseCond, types.NewStringDatum("case3")),
			result: types.Datum{},
		},
		{
			expr: buildExpr(tipb.ExprType_Case,
				buildExpr(tipb.ExprType_Case,
					falseCond, types.NewIntDatum(0),
					trueCond, types.NewIntDatum(1),
				), types.NewStringDatum("nested case when"),
				falseCond, types.NewStringDatum("case1"),
				trueCond, types.NewStringDatum("case2"),
				trueCond, types.NewStringDatum("case3")),
			result: types.NewStringDatum("nested case when"),
		},
		{
			expr: buildExpr(tipb.ExprType_Case,
				nullCond, types.NewStringDatum("case1"),
				falseCond, types.NewStringDatum("case2"),
				trueCond, types.NewStringDatum("case3")),
			result: types.NewStringDatum("case3"),
		},
	}
	for _, ca := range cases {
		result, err := xevaluator.Eval(ca.expr)
		c.Assert(err, IsNil)
		c.Assert(result.Kind(), Equals, ca.result.Kind())
		cmp, err := result.CompareDatum(ca.result)
		c.Assert(err, IsNil)
		c.Assert(cmp, Equals, 0)
	}
}

func (s *testEvalSuite) TestEvalCoalesce(c *C) {
	colID := int64(1)
	row := make(map[int64]types.Datum)
	row[colID] = types.NewIntDatum(100)
	xevaluator := &Evaluator{Row: row}
	nullDatum := types.Datum{}
	nullDatum.SetNull()
	notNullDatum := types.NewStringDatum("not-null")
	cases := []struct {
		expr   *tipb.Expr
		result types.Datum
	}{
		{
			expr:   buildExpr(tipb.ExprType_Coalesce, nullDatum, nullDatum, nullDatum),
			result: nullDatum,
		},
		{
			expr:   buildExpr(tipb.ExprType_Coalesce, nullDatum, notNullDatum, nullDatum),
			result: notNullDatum,
		},
		{
			expr:   buildExpr(tipb.ExprType_Coalesce, nullDatum, notNullDatum, types.NewStringDatum("not-null-2"), nullDatum),
			result: notNullDatum,
		},
	}
	for _, ca := range cases {
		result, err := xevaluator.Eval(ca.expr)
		c.Assert(err, IsNil)
		c.Assert(result.Kind(), Equals, ca.result.Kind())
		cmp, err := result.CompareDatum(ca.result)
		c.Assert(err, IsNil)
		c.Assert(cmp, Equals, 0)
	}
}
