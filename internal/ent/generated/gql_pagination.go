// Copyright 2023 The Infratographer Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by entc, DO NOT EDIT.

package generated

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.infratographer.com/ipam-api/internal/ent/generated/ipaddress"
	"go.infratographer.com/ipam-api/internal/ent/generated/ipblock"
	"go.infratographer.com/ipam-api/internal/ent/generated/ipblocktype"
	"go.infratographer.com/x/gidx"
)

// Common entgql types.
type (
	Cursor         = entgql.Cursor[gidx.PrefixedID]
	PageInfo       = entgql.PageInfo[gidx.PrefixedID]
	OrderDirection = entgql.OrderDirection
)

func orderFunc(o OrderDirection, field string) func(*sql.Selector) {
	if o == entgql.OrderDirectionDesc {
		return Desc(field)
	}
	return Asc(field)
}

const errInvalidPagination = "INVALID_PAGINATION"

func validateFirstLast(first, last *int) (err *gqlerror.Error) {
	switch {
	case first != nil && last != nil:
		err = &gqlerror.Error{
			Message: "Passing both `first` and `last` to paginate a connection is not supported.",
		}
	case first != nil && *first < 0:
		err = &gqlerror.Error{
			Message: "`first` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	case last != nil && *last < 0:
		err = &gqlerror.Error{
			Message: "`last` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	}
	return err
}

func collectedField(ctx context.Context, path ...string) *graphql.CollectedField {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return nil
	}
	field := fc.Field
	oc := graphql.GetOperationContext(ctx)
walk:
	for _, name := range path {
		for _, f := range graphql.CollectFields(oc, field.Selections, nil) {
			if f.Alias == name {
				field = f
				continue walk
			}
		}
		return nil
	}
	return &field
}

func hasCollectedField(ctx context.Context, path ...string) bool {
	if graphql.GetFieldContext(ctx) == nil {
		return true
	}
	return collectedField(ctx, path...) != nil
}

const (
	edgesField      = "edges"
	nodeField       = "node"
	pageInfoField   = "pageInfo"
	totalCountField = "totalCount"
)

func paginateLimit(first, last *int) int {
	var limit int
	if first != nil {
		limit = *first + 1
	} else if last != nil {
		limit = *last + 1
	}
	return limit
}

// IPAddressEdge is the edge representation of IPAddress.
type IPAddressEdge struct {
	Node   *IPAddress `json:"node"`
	Cursor Cursor     `json:"cursor"`
}

// IPAddressConnection is the connection containing edges to IPAddress.
type IPAddressConnection struct {
	Edges      []*IPAddressEdge `json:"edges"`
	PageInfo   PageInfo         `json:"pageInfo"`
	TotalCount int              `json:"totalCount"`
}

func (c *IPAddressConnection) build(nodes []*IPAddress, pager *ipaddressPager, after *Cursor, first *int, before *Cursor, last *int) {
	c.PageInfo.HasNextPage = before != nil
	c.PageInfo.HasPreviousPage = after != nil
	if first != nil && *first+1 == len(nodes) {
		c.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && *last+1 == len(nodes) {
		c.PageInfo.HasPreviousPage = true
		nodes = nodes[:len(nodes)-1]
	}
	var nodeAt func(int) *IPAddress
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *IPAddress {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *IPAddress {
			return nodes[i]
		}
	}
	c.Edges = make([]*IPAddressEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		c.Edges[i] = &IPAddressEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}
	if l := len(c.Edges); l > 0 {
		c.PageInfo.StartCursor = &c.Edges[0].Cursor
		c.PageInfo.EndCursor = &c.Edges[l-1].Cursor
	}
	if c.TotalCount == 0 {
		c.TotalCount = len(nodes)
	}
}

// IPAddressPaginateOption enables pagination customization.
type IPAddressPaginateOption func(*ipaddressPager) error

// WithIPAddressOrder configures pagination ordering.
func WithIPAddressOrder(order *IPAddressOrder) IPAddressPaginateOption {
	if order == nil {
		order = DefaultIPAddressOrder
	}
	o := *order
	return func(pager *ipaddressPager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultIPAddressOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithIPAddressFilter configures pagination filter.
func WithIPAddressFilter(filter func(*IPAddressQuery) (*IPAddressQuery, error)) IPAddressPaginateOption {
	return func(pager *ipaddressPager) error {
		if filter == nil {
			return errors.New("IPAddressQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type ipaddressPager struct {
	reverse bool
	order   *IPAddressOrder
	filter  func(*IPAddressQuery) (*IPAddressQuery, error)
}

func newIPAddressPager(opts []IPAddressPaginateOption, reverse bool) (*ipaddressPager, error) {
	pager := &ipaddressPager{reverse: reverse}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultIPAddressOrder
	}
	return pager, nil
}

func (p *ipaddressPager) applyFilter(query *IPAddressQuery) (*IPAddressQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *ipaddressPager) toCursor(ia *IPAddress) Cursor {
	return p.order.Field.toCursor(ia)
}

func (p *ipaddressPager) applyCursors(query *IPAddressQuery, after, before *Cursor) (*IPAddressQuery, error) {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	for _, predicate := range entgql.CursorsPredicate(after, before, DefaultIPAddressOrder.Field.column, p.order.Field.column, direction) {
		query = query.Where(predicate)
	}
	return query, nil
}

func (p *ipaddressPager) applyOrder(query *IPAddressQuery) *IPAddressQuery {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	query = query.Order(p.order.Field.toTerm(direction.OrderTermOption()))
	if p.order.Field != DefaultIPAddressOrder.Field {
		query = query.Order(DefaultIPAddressOrder.Field.toTerm(direction.OrderTermOption()))
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(p.order.Field.column)
	}
	return query
}

func (p *ipaddressPager) orderExpr(query *IPAddressQuery) sql.Querier {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(p.order.Field.column)
	}
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(p.order.Field.column).Pad().WriteString(string(direction))
		if p.order.Field != DefaultIPAddressOrder.Field {
			b.Comma().Ident(DefaultIPAddressOrder.Field.column).Pad().WriteString(string(direction))
		}
	})
}

// Paginate executes the query and returns a relay based cursor connection to IPAddress.
func (ia *IPAddressQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...IPAddressPaginateOption,
) (*IPAddressConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newIPAddressPager(opts, last != nil)
	if err != nil {
		return nil, err
	}
	if ia, err = pager.applyFilter(ia); err != nil {
		return nil, err
	}
	conn := &IPAddressConnection{Edges: []*IPAddressEdge{}}
	ignoredEdges := !hasCollectedField(ctx, edgesField)
	if hasCollectedField(ctx, totalCountField) || hasCollectedField(ctx, pageInfoField) {
		hasPagination := after != nil || first != nil || before != nil || last != nil
		if hasPagination || ignoredEdges {
			if conn.TotalCount, err = ia.Clone().Count(ctx); err != nil {
				return nil, err
			}
			conn.PageInfo.HasNextPage = first != nil && conn.TotalCount > 0
			conn.PageInfo.HasPreviousPage = last != nil && conn.TotalCount > 0
		}
	}
	if ignoredEdges || (first != nil && *first == 0) || (last != nil && *last == 0) {
		return conn, nil
	}
	if ia, err = pager.applyCursors(ia, after, before); err != nil {
		return nil, err
	}
	if limit := paginateLimit(first, last); limit != 0 {
		ia.Limit(limit)
	}
	if field := collectedField(ctx, edgesField, nodeField); field != nil {
		if err := ia.collectField(ctx, graphql.GetOperationContext(ctx), *field, []string{edgesField, nodeField}); err != nil {
			return nil, err
		}
	}
	ia = pager.applyOrder(ia)
	nodes, err := ia.All(ctx)
	if err != nil {
		return nil, err
	}
	conn.build(nodes, pager, after, first, before, last)
	return conn, nil
}

var (
	// IPAddressOrderFieldID orders IPAddress by id.
	IPAddressOrderFieldID = &IPAddressOrderField{
		Value: func(ia *IPAddress) (ent.Value, error) {
			return ia.ID, nil
		},
		column: ipaddress.FieldID,
		toTerm: ipaddress.ByID,
		toCursor: func(ia *IPAddress) Cursor {
			return Cursor{
				ID:    ia.ID,
				Value: ia.ID,
			}
		},
	}
	// IPAddressOrderFieldCreatedAt orders IPAddress by created_at.
	IPAddressOrderFieldCreatedAt = &IPAddressOrderField{
		Value: func(ia *IPAddress) (ent.Value, error) {
			return ia.CreatedAt, nil
		},
		column: ipaddress.FieldCreatedAt,
		toTerm: ipaddress.ByCreatedAt,
		toCursor: func(ia *IPAddress) Cursor {
			return Cursor{
				ID:    ia.ID,
				Value: ia.CreatedAt,
			}
		},
	}
	// IPAddressOrderFieldUpdatedAt orders IPAddress by updated_at.
	IPAddressOrderFieldUpdatedAt = &IPAddressOrderField{
		Value: func(ia *IPAddress) (ent.Value, error) {
			return ia.UpdatedAt, nil
		},
		column: ipaddress.FieldUpdatedAt,
		toTerm: ipaddress.ByUpdatedAt,
		toCursor: func(ia *IPAddress) Cursor {
			return Cursor{
				ID:    ia.ID,
				Value: ia.UpdatedAt,
			}
		},
	}
	// IPAddressOrderFieldIP orders IPAddress by IP.
	IPAddressOrderFieldIP = &IPAddressOrderField{
		Value: func(ia *IPAddress) (ent.Value, error) {
			return ia.IP, nil
		},
		column: ipaddress.FieldIP,
		toTerm: ipaddress.ByIP,
		toCursor: func(ia *IPAddress) Cursor {
			return Cursor{
				ID:    ia.ID,
				Value: ia.IP,
			}
		},
	}
	// IPAddressOrderFieldBlockID orders IPAddress by block_id.
	IPAddressOrderFieldBlockID = &IPAddressOrderField{
		Value: func(ia *IPAddress) (ent.Value, error) {
			return ia.BlockID, nil
		},
		column: ipaddress.FieldBlockID,
		toTerm: ipaddress.ByBlockID,
		toCursor: func(ia *IPAddress) Cursor {
			return Cursor{
				ID:    ia.ID,
				Value: ia.BlockID,
			}
		},
	}
	// IPAddressOrderFieldNodeID orders IPAddress by node_id.
	IPAddressOrderFieldNodeID = &IPAddressOrderField{
		Value: func(ia *IPAddress) (ent.Value, error) {
			return ia.NodeID, nil
		},
		column: ipaddress.FieldNodeID,
		toTerm: ipaddress.ByNodeID,
		toCursor: func(ia *IPAddress) Cursor {
			return Cursor{
				ID:    ia.ID,
				Value: ia.NodeID,
			}
		},
	}
	// IPAddressOrderFieldNodeOwnerID orders IPAddress by node_owner_id.
	IPAddressOrderFieldNodeOwnerID = &IPAddressOrderField{
		Value: func(ia *IPAddress) (ent.Value, error) {
			return ia.NodeOwnerID, nil
		},
		column: ipaddress.FieldNodeOwnerID,
		toTerm: ipaddress.ByNodeOwnerID,
		toCursor: func(ia *IPAddress) Cursor {
			return Cursor{
				ID:    ia.ID,
				Value: ia.NodeOwnerID,
			}
		},
	}
	// IPAddressOrderFieldReserved orders IPAddress by reserved.
	IPAddressOrderFieldReserved = &IPAddressOrderField{
		Value: func(ia *IPAddress) (ent.Value, error) {
			return ia.Reserved, nil
		},
		column: ipaddress.FieldReserved,
		toTerm: ipaddress.ByReserved,
		toCursor: func(ia *IPAddress) Cursor {
			return Cursor{
				ID:    ia.ID,
				Value: ia.Reserved,
			}
		},
	}
)

// String implement fmt.Stringer interface.
func (f IPAddressOrderField) String() string {
	var str string
	switch f.column {
	case IPAddressOrderFieldID.column:
		str = "ID"
	case IPAddressOrderFieldCreatedAt.column:
		str = "CREATED_AT"
	case IPAddressOrderFieldUpdatedAt.column:
		str = "UPDATED_AT"
	case IPAddressOrderFieldIP.column:
		str = "IP"
	case IPAddressOrderFieldBlockID.column:
		str = "BLOCK"
	case IPAddressOrderFieldNodeID.column:
		str = "NODE"
	case IPAddressOrderFieldNodeOwnerID.column:
		str = "OWNER"
	case IPAddressOrderFieldReserved.column:
		str = "RESERVED"
	}
	return str
}

// MarshalGQL implements graphql.Marshaler interface.
func (f IPAddressOrderField) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(f.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (f *IPAddressOrderField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("IPAddressOrderField %T must be a string", v)
	}
	switch str {
	case "ID":
		*f = *IPAddressOrderFieldID
	case "CREATED_AT":
		*f = *IPAddressOrderFieldCreatedAt
	case "UPDATED_AT":
		*f = *IPAddressOrderFieldUpdatedAt
	case "IP":
		*f = *IPAddressOrderFieldIP
	case "BLOCK":
		*f = *IPAddressOrderFieldBlockID
	case "NODE":
		*f = *IPAddressOrderFieldNodeID
	case "OWNER":
		*f = *IPAddressOrderFieldNodeOwnerID
	case "RESERVED":
		*f = *IPAddressOrderFieldReserved
	default:
		return fmt.Errorf("%s is not a valid IPAddressOrderField", str)
	}
	return nil
}

// IPAddressOrderField defines the ordering field of IPAddress.
type IPAddressOrderField struct {
	// Value extracts the ordering value from the given IPAddress.
	Value    func(*IPAddress) (ent.Value, error)
	column   string // field or computed.
	toTerm   func(...sql.OrderTermOption) ipaddress.OrderOption
	toCursor func(*IPAddress) Cursor
}

// IPAddressOrder defines the ordering of IPAddress.
type IPAddressOrder struct {
	Direction OrderDirection       `json:"direction"`
	Field     *IPAddressOrderField `json:"field"`
}

// DefaultIPAddressOrder is the default ordering of IPAddress.
var DefaultIPAddressOrder = &IPAddressOrder{
	Direction: entgql.OrderDirectionAsc,
	Field: &IPAddressOrderField{
		Value: func(ia *IPAddress) (ent.Value, error) {
			return ia.ID, nil
		},
		column: ipaddress.FieldID,
		toTerm: ipaddress.ByID,
		toCursor: func(ia *IPAddress) Cursor {
			return Cursor{ID: ia.ID}
		},
	},
}

// ToEdge converts IPAddress into IPAddressEdge.
func (ia *IPAddress) ToEdge(order *IPAddressOrder) *IPAddressEdge {
	if order == nil {
		order = DefaultIPAddressOrder
	}
	return &IPAddressEdge{
		Node:   ia,
		Cursor: order.Field.toCursor(ia),
	}
}

// IPBlockEdge is the edge representation of IPBlock.
type IPBlockEdge struct {
	Node   *IPBlock `json:"node"`
	Cursor Cursor   `json:"cursor"`
}

// IPBlockConnection is the connection containing edges to IPBlock.
type IPBlockConnection struct {
	Edges      []*IPBlockEdge `json:"edges"`
	PageInfo   PageInfo       `json:"pageInfo"`
	TotalCount int            `json:"totalCount"`
}

func (c *IPBlockConnection) build(nodes []*IPBlock, pager *ipblockPager, after *Cursor, first *int, before *Cursor, last *int) {
	c.PageInfo.HasNextPage = before != nil
	c.PageInfo.HasPreviousPage = after != nil
	if first != nil && *first+1 == len(nodes) {
		c.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && *last+1 == len(nodes) {
		c.PageInfo.HasPreviousPage = true
		nodes = nodes[:len(nodes)-1]
	}
	var nodeAt func(int) *IPBlock
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *IPBlock {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *IPBlock {
			return nodes[i]
		}
	}
	c.Edges = make([]*IPBlockEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		c.Edges[i] = &IPBlockEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}
	if l := len(c.Edges); l > 0 {
		c.PageInfo.StartCursor = &c.Edges[0].Cursor
		c.PageInfo.EndCursor = &c.Edges[l-1].Cursor
	}
	if c.TotalCount == 0 {
		c.TotalCount = len(nodes)
	}
}

// IPBlockPaginateOption enables pagination customization.
type IPBlockPaginateOption func(*ipblockPager) error

// WithIPBlockOrder configures pagination ordering.
func WithIPBlockOrder(order *IPBlockOrder) IPBlockPaginateOption {
	if order == nil {
		order = DefaultIPBlockOrder
	}
	o := *order
	return func(pager *ipblockPager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultIPBlockOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithIPBlockFilter configures pagination filter.
func WithIPBlockFilter(filter func(*IPBlockQuery) (*IPBlockQuery, error)) IPBlockPaginateOption {
	return func(pager *ipblockPager) error {
		if filter == nil {
			return errors.New("IPBlockQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type ipblockPager struct {
	reverse bool
	order   *IPBlockOrder
	filter  func(*IPBlockQuery) (*IPBlockQuery, error)
}

func newIPBlockPager(opts []IPBlockPaginateOption, reverse bool) (*ipblockPager, error) {
	pager := &ipblockPager{reverse: reverse}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultIPBlockOrder
	}
	return pager, nil
}

func (p *ipblockPager) applyFilter(query *IPBlockQuery) (*IPBlockQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *ipblockPager) toCursor(ib *IPBlock) Cursor {
	return p.order.Field.toCursor(ib)
}

func (p *ipblockPager) applyCursors(query *IPBlockQuery, after, before *Cursor) (*IPBlockQuery, error) {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	for _, predicate := range entgql.CursorsPredicate(after, before, DefaultIPBlockOrder.Field.column, p.order.Field.column, direction) {
		query = query.Where(predicate)
	}
	return query, nil
}

func (p *ipblockPager) applyOrder(query *IPBlockQuery) *IPBlockQuery {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	query = query.Order(p.order.Field.toTerm(direction.OrderTermOption()))
	if p.order.Field != DefaultIPBlockOrder.Field {
		query = query.Order(DefaultIPBlockOrder.Field.toTerm(direction.OrderTermOption()))
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(p.order.Field.column)
	}
	return query
}

func (p *ipblockPager) orderExpr(query *IPBlockQuery) sql.Querier {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(p.order.Field.column)
	}
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(p.order.Field.column).Pad().WriteString(string(direction))
		if p.order.Field != DefaultIPBlockOrder.Field {
			b.Comma().Ident(DefaultIPBlockOrder.Field.column).Pad().WriteString(string(direction))
		}
	})
}

// Paginate executes the query and returns a relay based cursor connection to IPBlock.
func (ib *IPBlockQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...IPBlockPaginateOption,
) (*IPBlockConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newIPBlockPager(opts, last != nil)
	if err != nil {
		return nil, err
	}
	if ib, err = pager.applyFilter(ib); err != nil {
		return nil, err
	}
	conn := &IPBlockConnection{Edges: []*IPBlockEdge{}}
	ignoredEdges := !hasCollectedField(ctx, edgesField)
	if hasCollectedField(ctx, totalCountField) || hasCollectedField(ctx, pageInfoField) {
		hasPagination := after != nil || first != nil || before != nil || last != nil
		if hasPagination || ignoredEdges {
			if conn.TotalCount, err = ib.Clone().Count(ctx); err != nil {
				return nil, err
			}
			conn.PageInfo.HasNextPage = first != nil && conn.TotalCount > 0
			conn.PageInfo.HasPreviousPage = last != nil && conn.TotalCount > 0
		}
	}
	if ignoredEdges || (first != nil && *first == 0) || (last != nil && *last == 0) {
		return conn, nil
	}
	if ib, err = pager.applyCursors(ib, after, before); err != nil {
		return nil, err
	}
	if limit := paginateLimit(first, last); limit != 0 {
		ib.Limit(limit)
	}
	if field := collectedField(ctx, edgesField, nodeField); field != nil {
		if err := ib.collectField(ctx, graphql.GetOperationContext(ctx), *field, []string{edgesField, nodeField}); err != nil {
			return nil, err
		}
	}
	ib = pager.applyOrder(ib)
	nodes, err := ib.All(ctx)
	if err != nil {
		return nil, err
	}
	conn.build(nodes, pager, after, first, before, last)
	return conn, nil
}

var (
	// IPBlockOrderFieldID orders IPBlock by id.
	IPBlockOrderFieldID = &IPBlockOrderField{
		Value: func(ib *IPBlock) (ent.Value, error) {
			return ib.ID, nil
		},
		column: ipblock.FieldID,
		toTerm: ipblock.ByID,
		toCursor: func(ib *IPBlock) Cursor {
			return Cursor{
				ID:    ib.ID,
				Value: ib.ID,
			}
		},
	}
	// IPBlockOrderFieldCreatedAt orders IPBlock by created_at.
	IPBlockOrderFieldCreatedAt = &IPBlockOrderField{
		Value: func(ib *IPBlock) (ent.Value, error) {
			return ib.CreatedAt, nil
		},
		column: ipblock.FieldCreatedAt,
		toTerm: ipblock.ByCreatedAt,
		toCursor: func(ib *IPBlock) Cursor {
			return Cursor{
				ID:    ib.ID,
				Value: ib.CreatedAt,
			}
		},
	}
	// IPBlockOrderFieldUpdatedAt orders IPBlock by updated_at.
	IPBlockOrderFieldUpdatedAt = &IPBlockOrderField{
		Value: func(ib *IPBlock) (ent.Value, error) {
			return ib.UpdatedAt, nil
		},
		column: ipblock.FieldUpdatedAt,
		toTerm: ipblock.ByUpdatedAt,
		toCursor: func(ib *IPBlock) Cursor {
			return Cursor{
				ID:    ib.ID,
				Value: ib.UpdatedAt,
			}
		},
	}
	// IPBlockOrderFieldPrefix orders IPBlock by prefix.
	IPBlockOrderFieldPrefix = &IPBlockOrderField{
		Value: func(ib *IPBlock) (ent.Value, error) {
			return ib.Prefix, nil
		},
		column: ipblock.FieldPrefix,
		toTerm: ipblock.ByPrefix,
		toCursor: func(ib *IPBlock) Cursor {
			return Cursor{
				ID:    ib.ID,
				Value: ib.Prefix,
			}
		},
	}
	// IPBlockOrderFieldBlockTypeID orders IPBlock by block_type_id.
	IPBlockOrderFieldBlockTypeID = &IPBlockOrderField{
		Value: func(ib *IPBlock) (ent.Value, error) {
			return ib.BlockTypeID, nil
		},
		column: ipblock.FieldBlockTypeID,
		toTerm: ipblock.ByBlockTypeID,
		toCursor: func(ib *IPBlock) Cursor {
			return Cursor{
				ID:    ib.ID,
				Value: ib.BlockTypeID,
			}
		},
	}
	// IPBlockOrderFieldLocationID orders IPBlock by location_id.
	IPBlockOrderFieldLocationID = &IPBlockOrderField{
		Value: func(ib *IPBlock) (ent.Value, error) {
			return ib.LocationID, nil
		},
		column: ipblock.FieldLocationID,
		toTerm: ipblock.ByLocationID,
		toCursor: func(ib *IPBlock) Cursor {
			return Cursor{
				ID:    ib.ID,
				Value: ib.LocationID,
			}
		},
	}
	// IPBlockOrderFieldParentBlockID orders IPBlock by parent_block_id.
	IPBlockOrderFieldParentBlockID = &IPBlockOrderField{
		Value: func(ib *IPBlock) (ent.Value, error) {
			return ib.ParentBlockID, nil
		},
		column: ipblock.FieldParentBlockID,
		toTerm: ipblock.ByParentBlockID,
		toCursor: func(ib *IPBlock) Cursor {
			return Cursor{
				ID:    ib.ID,
				Value: ib.ParentBlockID,
			}
		},
	}
	// IPBlockOrderFieldAllowAutoSubnet orders IPBlock by allow_auto_subnet.
	IPBlockOrderFieldAllowAutoSubnet = &IPBlockOrderField{
		Value: func(ib *IPBlock) (ent.Value, error) {
			return ib.AllowAutoSubnet, nil
		},
		column: ipblock.FieldAllowAutoSubnet,
		toTerm: ipblock.ByAllowAutoSubnet,
		toCursor: func(ib *IPBlock) Cursor {
			return Cursor{
				ID:    ib.ID,
				Value: ib.AllowAutoSubnet,
			}
		},
	}
	// IPBlockOrderFieldAllowAutoAllocate orders IPBlock by allow_auto_allocate.
	IPBlockOrderFieldAllowAutoAllocate = &IPBlockOrderField{
		Value: func(ib *IPBlock) (ent.Value, error) {
			return ib.AllowAutoAllocate, nil
		},
		column: ipblock.FieldAllowAutoAllocate,
		toTerm: ipblock.ByAllowAutoAllocate,
		toCursor: func(ib *IPBlock) Cursor {
			return Cursor{
				ID:    ib.ID,
				Value: ib.AllowAutoAllocate,
			}
		},
	}
)

// String implement fmt.Stringer interface.
func (f IPBlockOrderField) String() string {
	var str string
	switch f.column {
	case IPBlockOrderFieldID.column:
		str = "ID"
	case IPBlockOrderFieldCreatedAt.column:
		str = "CREATED_AT"
	case IPBlockOrderFieldUpdatedAt.column:
		str = "UPDATED_AT"
	case IPBlockOrderFieldPrefix.column:
		str = "PREFIX"
	case IPBlockOrderFieldBlockTypeID.column:
		str = "BLOCK_TYPE"
	case IPBlockOrderFieldLocationID.column:
		str = "LOCATION"
	case IPBlockOrderFieldParentBlockID.column:
		str = "PARENT_BLOCK"
	case IPBlockOrderFieldAllowAutoSubnet.column:
		str = "AUTOSUBNET"
	case IPBlockOrderFieldAllowAutoAllocate.column:
		str = "AUTOALLOCATE"
	}
	return str
}

// MarshalGQL implements graphql.Marshaler interface.
func (f IPBlockOrderField) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(f.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (f *IPBlockOrderField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("IPBlockOrderField %T must be a string", v)
	}
	switch str {
	case "ID":
		*f = *IPBlockOrderFieldID
	case "CREATED_AT":
		*f = *IPBlockOrderFieldCreatedAt
	case "UPDATED_AT":
		*f = *IPBlockOrderFieldUpdatedAt
	case "PREFIX":
		*f = *IPBlockOrderFieldPrefix
	case "BLOCK_TYPE":
		*f = *IPBlockOrderFieldBlockTypeID
	case "LOCATION":
		*f = *IPBlockOrderFieldLocationID
	case "PARENT_BLOCK":
		*f = *IPBlockOrderFieldParentBlockID
	case "AUTOSUBNET":
		*f = *IPBlockOrderFieldAllowAutoSubnet
	case "AUTOALLOCATE":
		*f = *IPBlockOrderFieldAllowAutoAllocate
	default:
		return fmt.Errorf("%s is not a valid IPBlockOrderField", str)
	}
	return nil
}

// IPBlockOrderField defines the ordering field of IPBlock.
type IPBlockOrderField struct {
	// Value extracts the ordering value from the given IPBlock.
	Value    func(*IPBlock) (ent.Value, error)
	column   string // field or computed.
	toTerm   func(...sql.OrderTermOption) ipblock.OrderOption
	toCursor func(*IPBlock) Cursor
}

// IPBlockOrder defines the ordering of IPBlock.
type IPBlockOrder struct {
	Direction OrderDirection     `json:"direction"`
	Field     *IPBlockOrderField `json:"field"`
}

// DefaultIPBlockOrder is the default ordering of IPBlock.
var DefaultIPBlockOrder = &IPBlockOrder{
	Direction: entgql.OrderDirectionAsc,
	Field: &IPBlockOrderField{
		Value: func(ib *IPBlock) (ent.Value, error) {
			return ib.ID, nil
		},
		column: ipblock.FieldID,
		toTerm: ipblock.ByID,
		toCursor: func(ib *IPBlock) Cursor {
			return Cursor{ID: ib.ID}
		},
	},
}

// ToEdge converts IPBlock into IPBlockEdge.
func (ib *IPBlock) ToEdge(order *IPBlockOrder) *IPBlockEdge {
	if order == nil {
		order = DefaultIPBlockOrder
	}
	return &IPBlockEdge{
		Node:   ib,
		Cursor: order.Field.toCursor(ib),
	}
}

// IPBlockTypeEdge is the edge representation of IPBlockType.
type IPBlockTypeEdge struct {
	Node   *IPBlockType `json:"node"`
	Cursor Cursor       `json:"cursor"`
}

// IPBlockTypeConnection is the connection containing edges to IPBlockType.
type IPBlockTypeConnection struct {
	Edges      []*IPBlockTypeEdge `json:"edges"`
	PageInfo   PageInfo           `json:"pageInfo"`
	TotalCount int                `json:"totalCount"`
}

func (c *IPBlockTypeConnection) build(nodes []*IPBlockType, pager *ipblocktypePager, after *Cursor, first *int, before *Cursor, last *int) {
	c.PageInfo.HasNextPage = before != nil
	c.PageInfo.HasPreviousPage = after != nil
	if first != nil && *first+1 == len(nodes) {
		c.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && *last+1 == len(nodes) {
		c.PageInfo.HasPreviousPage = true
		nodes = nodes[:len(nodes)-1]
	}
	var nodeAt func(int) *IPBlockType
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *IPBlockType {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *IPBlockType {
			return nodes[i]
		}
	}
	c.Edges = make([]*IPBlockTypeEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		c.Edges[i] = &IPBlockTypeEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}
	if l := len(c.Edges); l > 0 {
		c.PageInfo.StartCursor = &c.Edges[0].Cursor
		c.PageInfo.EndCursor = &c.Edges[l-1].Cursor
	}
	if c.TotalCount == 0 {
		c.TotalCount = len(nodes)
	}
}

// IPBlockTypePaginateOption enables pagination customization.
type IPBlockTypePaginateOption func(*ipblocktypePager) error

// WithIPBlockTypeOrder configures pagination ordering.
func WithIPBlockTypeOrder(order *IPBlockTypeOrder) IPBlockTypePaginateOption {
	if order == nil {
		order = DefaultIPBlockTypeOrder
	}
	o := *order
	return func(pager *ipblocktypePager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultIPBlockTypeOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithIPBlockTypeFilter configures pagination filter.
func WithIPBlockTypeFilter(filter func(*IPBlockTypeQuery) (*IPBlockTypeQuery, error)) IPBlockTypePaginateOption {
	return func(pager *ipblocktypePager) error {
		if filter == nil {
			return errors.New("IPBlockTypeQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type ipblocktypePager struct {
	reverse bool
	order   *IPBlockTypeOrder
	filter  func(*IPBlockTypeQuery) (*IPBlockTypeQuery, error)
}

func newIPBlockTypePager(opts []IPBlockTypePaginateOption, reverse bool) (*ipblocktypePager, error) {
	pager := &ipblocktypePager{reverse: reverse}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultIPBlockTypeOrder
	}
	return pager, nil
}

func (p *ipblocktypePager) applyFilter(query *IPBlockTypeQuery) (*IPBlockTypeQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *ipblocktypePager) toCursor(ibt *IPBlockType) Cursor {
	return p.order.Field.toCursor(ibt)
}

func (p *ipblocktypePager) applyCursors(query *IPBlockTypeQuery, after, before *Cursor) (*IPBlockTypeQuery, error) {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	for _, predicate := range entgql.CursorsPredicate(after, before, DefaultIPBlockTypeOrder.Field.column, p.order.Field.column, direction) {
		query = query.Where(predicate)
	}
	return query, nil
}

func (p *ipblocktypePager) applyOrder(query *IPBlockTypeQuery) *IPBlockTypeQuery {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	query = query.Order(p.order.Field.toTerm(direction.OrderTermOption()))
	if p.order.Field != DefaultIPBlockTypeOrder.Field {
		query = query.Order(DefaultIPBlockTypeOrder.Field.toTerm(direction.OrderTermOption()))
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(p.order.Field.column)
	}
	return query
}

func (p *ipblocktypePager) orderExpr(query *IPBlockTypeQuery) sql.Querier {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(p.order.Field.column)
	}
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(p.order.Field.column).Pad().WriteString(string(direction))
		if p.order.Field != DefaultIPBlockTypeOrder.Field {
			b.Comma().Ident(DefaultIPBlockTypeOrder.Field.column).Pad().WriteString(string(direction))
		}
	})
}

// Paginate executes the query and returns a relay based cursor connection to IPBlockType.
func (ibt *IPBlockTypeQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...IPBlockTypePaginateOption,
) (*IPBlockTypeConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newIPBlockTypePager(opts, last != nil)
	if err != nil {
		return nil, err
	}
	if ibt, err = pager.applyFilter(ibt); err != nil {
		return nil, err
	}
	conn := &IPBlockTypeConnection{Edges: []*IPBlockTypeEdge{}}
	ignoredEdges := !hasCollectedField(ctx, edgesField)
	if hasCollectedField(ctx, totalCountField) || hasCollectedField(ctx, pageInfoField) {
		hasPagination := after != nil || first != nil || before != nil || last != nil
		if hasPagination || ignoredEdges {
			if conn.TotalCount, err = ibt.Clone().Count(ctx); err != nil {
				return nil, err
			}
			conn.PageInfo.HasNextPage = first != nil && conn.TotalCount > 0
			conn.PageInfo.HasPreviousPage = last != nil && conn.TotalCount > 0
		}
	}
	if ignoredEdges || (first != nil && *first == 0) || (last != nil && *last == 0) {
		return conn, nil
	}
	if ibt, err = pager.applyCursors(ibt, after, before); err != nil {
		return nil, err
	}
	if limit := paginateLimit(first, last); limit != 0 {
		ibt.Limit(limit)
	}
	if field := collectedField(ctx, edgesField, nodeField); field != nil {
		if err := ibt.collectField(ctx, graphql.GetOperationContext(ctx), *field, []string{edgesField, nodeField}); err != nil {
			return nil, err
		}
	}
	ibt = pager.applyOrder(ibt)
	nodes, err := ibt.All(ctx)
	if err != nil {
		return nil, err
	}
	conn.build(nodes, pager, after, first, before, last)
	return conn, nil
}

var (
	// IPBlockTypeOrderFieldID orders IPBlockType by id.
	IPBlockTypeOrderFieldID = &IPBlockTypeOrderField{
		Value: func(ibt *IPBlockType) (ent.Value, error) {
			return ibt.ID, nil
		},
		column: ipblocktype.FieldID,
		toTerm: ipblocktype.ByID,
		toCursor: func(ibt *IPBlockType) Cursor {
			return Cursor{
				ID:    ibt.ID,
				Value: ibt.ID,
			}
		},
	}
	// IPBlockTypeOrderFieldCreatedAt orders IPBlockType by created_at.
	IPBlockTypeOrderFieldCreatedAt = &IPBlockTypeOrderField{
		Value: func(ibt *IPBlockType) (ent.Value, error) {
			return ibt.CreatedAt, nil
		},
		column: ipblocktype.FieldCreatedAt,
		toTerm: ipblocktype.ByCreatedAt,
		toCursor: func(ibt *IPBlockType) Cursor {
			return Cursor{
				ID:    ibt.ID,
				Value: ibt.CreatedAt,
			}
		},
	}
	// IPBlockTypeOrderFieldUpdatedAt orders IPBlockType by updated_at.
	IPBlockTypeOrderFieldUpdatedAt = &IPBlockTypeOrderField{
		Value: func(ibt *IPBlockType) (ent.Value, error) {
			return ibt.UpdatedAt, nil
		},
		column: ipblocktype.FieldUpdatedAt,
		toTerm: ipblocktype.ByUpdatedAt,
		toCursor: func(ibt *IPBlockType) Cursor {
			return Cursor{
				ID:    ibt.ID,
				Value: ibt.UpdatedAt,
			}
		},
	}
	// IPBlockTypeOrderFieldName orders IPBlockType by name.
	IPBlockTypeOrderFieldName = &IPBlockTypeOrderField{
		Value: func(ibt *IPBlockType) (ent.Value, error) {
			return ibt.Name, nil
		},
		column: ipblocktype.FieldName,
		toTerm: ipblocktype.ByName,
		toCursor: func(ibt *IPBlockType) Cursor {
			return Cursor{
				ID:    ibt.ID,
				Value: ibt.Name,
			}
		},
	}
	// IPBlockTypeOrderFieldOwnerID orders IPBlockType by owner_id.
	IPBlockTypeOrderFieldOwnerID = &IPBlockTypeOrderField{
		Value: func(ibt *IPBlockType) (ent.Value, error) {
			return ibt.OwnerID, nil
		},
		column: ipblocktype.FieldOwnerID,
		toTerm: ipblocktype.ByOwnerID,
		toCursor: func(ibt *IPBlockType) Cursor {
			return Cursor{
				ID:    ibt.ID,
				Value: ibt.OwnerID,
			}
		},
	}
)

// String implement fmt.Stringer interface.
func (f IPBlockTypeOrderField) String() string {
	var str string
	switch f.column {
	case IPBlockTypeOrderFieldID.column:
		str = "ID"
	case IPBlockTypeOrderFieldCreatedAt.column:
		str = "CREATED_AT"
	case IPBlockTypeOrderFieldUpdatedAt.column:
		str = "UPDATED_AT"
	case IPBlockTypeOrderFieldName.column:
		str = "NAME"
	case IPBlockTypeOrderFieldOwnerID.column:
		str = "OWNER"
	}
	return str
}

// MarshalGQL implements graphql.Marshaler interface.
func (f IPBlockTypeOrderField) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(f.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (f *IPBlockTypeOrderField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("IPBlockTypeOrderField %T must be a string", v)
	}
	switch str {
	case "ID":
		*f = *IPBlockTypeOrderFieldID
	case "CREATED_AT":
		*f = *IPBlockTypeOrderFieldCreatedAt
	case "UPDATED_AT":
		*f = *IPBlockTypeOrderFieldUpdatedAt
	case "NAME":
		*f = *IPBlockTypeOrderFieldName
	case "OWNER":
		*f = *IPBlockTypeOrderFieldOwnerID
	default:
		return fmt.Errorf("%s is not a valid IPBlockTypeOrderField", str)
	}
	return nil
}

// IPBlockTypeOrderField defines the ordering field of IPBlockType.
type IPBlockTypeOrderField struct {
	// Value extracts the ordering value from the given IPBlockType.
	Value    func(*IPBlockType) (ent.Value, error)
	column   string // field or computed.
	toTerm   func(...sql.OrderTermOption) ipblocktype.OrderOption
	toCursor func(*IPBlockType) Cursor
}

// IPBlockTypeOrder defines the ordering of IPBlockType.
type IPBlockTypeOrder struct {
	Direction OrderDirection         `json:"direction"`
	Field     *IPBlockTypeOrderField `json:"field"`
}

// DefaultIPBlockTypeOrder is the default ordering of IPBlockType.
var DefaultIPBlockTypeOrder = &IPBlockTypeOrder{
	Direction: entgql.OrderDirectionAsc,
	Field: &IPBlockTypeOrderField{
		Value: func(ibt *IPBlockType) (ent.Value, error) {
			return ibt.ID, nil
		},
		column: ipblocktype.FieldID,
		toTerm: ipblocktype.ByID,
		toCursor: func(ibt *IPBlockType) Cursor {
			return Cursor{ID: ibt.ID}
		},
	},
}

// ToEdge converts IPBlockType into IPBlockTypeEdge.
func (ibt *IPBlockType) ToEdge(order *IPBlockTypeOrder) *IPBlockTypeEdge {
	if order == nil {
		order = DefaultIPBlockTypeOrder
	}
	return &IPBlockTypeEdge{
		Node:   ibt,
		Cursor: order.Field.toCursor(ibt),
	}
}
