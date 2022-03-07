package token

import (
	"fmt"
	"reflect"
)

// Some utility functions for manipulating the semantic stack, used in the
// semantic actions generated by tool.go

// Final AST constructs
const (
	FINAL_PROG                        Kind = "Prog"
	FINAL_STRUCT_OR_IMPL_OR_FUNC_LIST Kind = "StructOrImplOrFuncList"
	FINAL_FUNC_DEF                    Kind = "FuncDef"
	FINAL_STRUCT_DECL                 Kind = "StructDecl"
	FINAL_IMPL_DEF                    Kind = "ImplDef"

	FINAL_TYPE    Kind = "Type"
	FINAL_ID      Kind = "Id"
	FINAL_VOID    Kind = "Void"
	FINAL_INTEGER Kind = "Integer"
	FINAL_FLOAT   Kind = "Float"

	FINAL_INTNUM   Kind = "IntNum"
	FINAL_FLOATNUM Kind = "FloatNum"

	FINAL_FUNC_DEF_PARAM     Kind = "Param"
	FINAL_FUNC_DEF_PARAMLIST Kind = "ParamList"

	FINAL_STAT_BLOCK Kind = "StatBlock"

	// The difference between FINAL_FUNC_BODY and FINAL_STAT_BLOCK is that
	// function bodies can contain both variable declaration as well as
	// statements whereas statement blocks may only contain statements
	FINAL_FUNC_BODY Kind = "Body"

	FINAL_VAR_DECL  Kind = "VarDecl"
	FINAL_FUNC_DECL Kind = "FuncDecl"
	FINAL_STATEMENT Kind = "Statement"

	FINAL_WRITE  Kind = "Write"
	FINAL_IF     Kind = "If"
	FINAL_WHILE  Kind = "While"
	FINAL_READ   Kind = "Read"
	FINAL_RETURN Kind = "Return"

	FINAL_ASSIGN Kind = "Assign(=)"

	FINAL_EXPR       Kind = "Expr"
	FINAL_ARITH_EXPR Kind = "ArithExpr"
	FINAL_FACTOR     Kind = "Factor"
	FINAL_TERM       Kind = "Term"

	// ADDOP
	FINAL_PLUS  Kind = "Plus(+)"
	FINAL_MINUS Kind = "Minus(-)"
	FINAL_OR    Kind = "Or"

	// MULTOP
	FINAL_MULT Kind = "Mult(*)"
	FINAL_DIV  Kind = "Div(/)"
	FINAL_AND  Kind = "And(&)"

	FINAL_EQ  Kind = "Eq(==)"
	FINAL_NEQ Kind = "Neq(<>)"
	FINAL_LT  Kind = "Lt(<)"
	FINAL_GT  Kind = "Gt(>)"
	FINAL_LEQ Kind = "Leq(<=)"
	FINAL_GEQ Kind = "Geq(>=)"

	FINAL_RETURNTYPE Kind = "ReturnType"

	FINAL_NOT      Kind = "Not(!)"
	FINAL_NEGATIVE Kind = "SignNegative(-)"
	FINAL_POSITIVE Kind = "SignPositive(+)"

	FINAL_REL_EXPR Kind = "RelExpr"

	FINAL_DIM     Kind = "Dim"
	FINAL_DIMLIST Kind = "DimList"

	FINAL_INDEX     Kind = "Index"
	FINAL_INDEXLIST Kind = "IndexList"

	FINAL_SUBJECT   Kind = "Subject"
	FINAL_VARIABLE  Kind = "Variable"
	FINAL_FUNC_CALL Kind = "FuncCall"

	FINAL_FUNC_CALL_PARAM     Kind = "FuncCallParam"
	FINAL_FUNC_CALL_PARAMLIST Kind = "FuncCallParamList"

	FINAL_STATBLOCK Kind = "StatBlock"

	// Function definition list for ImplDef
	FINAL_FUNC_DEF_LIST Kind = "FuncDefList"

	FINAL_INHERITS Kind = "Inherits"
	FINAL_MEMBER   Kind = "Member"
	FINAL_MEMBERS  Kind = "Members"
	FINAL_PRIVATE  Kind = "Private"
	FINAL_PUBLIC   Kind = "Public"
)

// Use operandTypes to describe what the top of the stack should look like. Each
// element of operantTypes should be either a Kind or a []Kind
func pullUp(stack *[]*ASTNode, which int, operandTypes ...interface{}) {
	eat := len(operandTypes)
	lengthOrPanic(stack, eat)
	if which < 0 || which > eat-1 {
		panic(fmt.Errorf("which should be between %v and %v, but was %v", 0, eat-1, which))
	}

	topn := topN(stack, eat)

	// Type check
	for i, operand := range topn {
		operandType := operandTypes[i]
		switch opType := operandType.(type) {
		case Kind:
			typeCheck(operand, opType)
		case []Kind:
			typeCheck(operand, opType...)
		default:
			panic(fmt.Errorf("operandTypes should be either Kind or []Kind"))
		}
	}

	// We modify the pulled element
	pulled := topn[which]
	pulled.Children = make([]*ASTNode, 0, eat-1)
	pulled.Children = append(pulled.Children, topn[:which]...)
	if which < eat-1 {
		pulled.Children = append(pulled.Children, topn[which+1:]...)
	}

	// Trim stack and replace
	transform(stack, eat, pulled)
}

// Swaps only the type of the top node
func swapTop(stack *[]*ASTNode, newTopType Kind, acceptableTopTypes ...Kind) {
	l := lengthOrPanic(stack, 1)
	typeCheck(top(stack), acceptableTopTypes...)
	(*stack)[l-1].Type = newTopType
}

// Wraps the node at the top of the stack
func wrapTop(stack *[]*ASTNode, typee Kind, acceptableTopTypes ...Kind) {
	l := lengthOrPanic(stack, 1)
	top := top(stack)
	typeCheck(top, acceptableTopTypes...)
	(*stack)[l-1] = &ASTNode{
		Type:     typee,
		Children: []*ASTNode{top},
	}
}

// Pushes a leaf node onto the stack
func pushNode(stack *[]*ASTNode, typee Kind, tok Token) {
	*stack = append(*stack, &ASTNode{
		Type:  typee,
		Token: tok,
	})
}

// Merges the top two nodes into a list, or creates a list if there is not
// already a list
func list(stack *[]*ASTNode, listType Kind, elementTypes ...Kind) {
	l := lengthOrPanic(stack, 1)

	t := top(stack)
	if typeCheckNoPanic(t, elementTypes...) != nil {
		if typeCheckNoPanic(t, listType) != nil {
			// Then we need to place an empty list on top
			*stack = append(*stack, &ASTNode{
				Type:     listType,
				Children: []*ASTNode{},
			})
		}
		return
	}

	// Cut top
	*stack = (*stack)[:l-1]

	if l >= 2 {
		tt := top(stack)
		if tt.Type == listType {
			tt.Children = append(tt.Children, t)
			return
		}
	}

	// Otherwise, make a new list
	*stack = append(*stack, &ASTNode{
		Type:     listType,
		Children: []*ASTNode{t},
	})
}

// The same as trimAndSet except it doesn't append any children to the node you
// provide
func transform(stack *[]*ASTNode, n int, newTop *ASTNode) {
	l := len(*stack)
	(*stack)[l-n] = newTop
	*stack = (*stack)[:l-n+1]
}

// Consumes the top n nodes and replaces them with a single node
func transformWithChildren(stack *[]*ASTNode, n int, newTop *ASTNode) {
	newTop.Children = append(newTop.Children, (*stack)[len(*stack)-n:]...)
	transform(stack, n, newTop)
}

func typeCheckNoPanic(record *ASTNode, acceptableTypes ...Kind) error {
	acceptable := false
	for _, t := range acceptableTypes {
		if record.Type == t {
			acceptable = true
		}
	}
	if !acceptable {
		return fmt.Errorf(
			"record should be of type %v but got %v",
			acceptableTypes, record.Type)
	}
	return nil
}

func typeCheck(record *ASTNode, acceptableTypes ...Kind) {
	if err := typeCheckNoPanic(record, acceptableTypes...); err != nil {
		panic(err)
	}
}

func top(stack *[]*ASTNode) *ASTNode {
	return (*stack)[len(*stack)-1]
}

func topN(stack *[]*ASTNode, n int) []*ASTNode {
	return (*stack)[len(*stack)-n:]
}

func lengthOrPanic(stack *[]*ASTNode, desiredLength int) int {
	l := len(*stack)
	if l < desiredLength {
		panic(fmt.Errorf(
			"expected stack to be %v elements, but was %v. Here is the stack: %v",
			desiredLength, l, *stack))
	}
	return l
}

// Makes a subtree with pushType at the root, type checks the stack top using
// the provided type check kinds. eat == len(typeChecks) or the function will
// panic
func eat(stack *[]*ASTNode, eat int, pushType Kind, typeChecks ...interface{}) {
	l := len(typeChecks)
	if eat != l {
		panic(fmt.Errorf(""+
			"'eat' must be the same as len(typeChecks), "+
			"eat = %v and len(typeChecks) = %v",
			eat, l))
	}

	lengthOrPanic(stack, eat)
	topn := topN(stack, eat)
	for i, t := range typeChecks {
		switch tt := t.(type) {
		case []Kind:
			typeCheck(topn[i], tt...)
		case Kind:
			typeCheck(topn[i], tt)
		default:
			panic(fmt.Errorf(
				"typeChecks should be either Kind or []Kind, but got %v",
				reflect.TypeOf(t)))
		}
	}
	transformWithChildren(stack, eat, &ASTNode{Type: pushType})
}

var exprOperandTypes = []Kind{
	FINAL_TERM,
	FINAL_FACTOR,

	FINAL_VARIABLE,
	FINAL_FUNC_CALL,

	// ADDOP
	FINAL_PLUS,
	FINAL_MINUS,
	FINAL_OR,

	// MULTOP
	FINAL_MULT,
	FINAL_DIV,
	FINAL_AND,

	// LITERALS
	FINAL_INTNUM,
	FINAL_FLOATNUM,
}

var relOpTypes = []Kind{
	FINAL_EQ,
	FINAL_NEQ,
	FINAL_LT,
	FINAL_GT,
	FINAL_LEQ,
	FINAL_GEQ,
}

var statementTypes = []Kind{
	FINAL_IF,
	FINAL_WHILE,
	FINAL_READ,
	FINAL_WRITE,
	FINAL_RETURN,
	FINAL_ASSIGN,
	FINAL_FUNC_CALL,
}

// Use this map to override the default functionality in the SEM_DISPATCH map
var semDisptachOverride = map[Kind]SemanticAction{

	// Pushes an empty Inherits on the stack. Use this to avoid listMerging a
	// struct Id with no inherits list
	SEM_INHERITS_FRESH: func(stack *[]*ASTNode, tok Token) {
		pushNode(stack, FINAL_INHERITS, Token{})
	},

	// Variable declaration (inside struct)
	//
	// Consumes: Id, Type, DimList
	// Produces: VarDecl
	SEM_VAR_DECL_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		eat(stack, 3, FINAL_VAR_DECL, FINAL_ID, FINAL_TYPE, FINAL_DIMLIST)
	},

	// Variable declaration (inside struct)
	//
	// Consumes: Id, ParamList, ReturnType
	// Produces: FuncDecl
	SEM_FUNC_DECL_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		eat(stack, 3, FINAL_FUNC_DECL, FINAL_ID, FINAL_FUNC_DEF_PARAMLIST, FINAL_RETURNTYPE)
	},

	// Structure declaration
	//
	// Consumes: Id, Inherits, Members
	// Produces: Struct
	SEM_STRUCT_DECL_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		eat(stack, 3, FINAL_STRUCT_DECL, FINAL_ID, FINAL_INHERITS, FINAL_MEMBERS)
	},

	// Inheritance list for structs
	//
	// Consumes: (Inherits, Id) or (Id) or ()
	// Produces: Inherits
	SEM_INHERITS_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		list(stack, FINAL_INHERITS, FINAL_ID)
	},

	// Members list for structs
	//
	// Consumes: (Members, Member) or (Member) or ()
	// Produces: Members
	SEM_MEMBERS_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		list(stack, FINAL_MEMBERS, FINAL_MEMBER)
	},

	// Member of a struct
	//
	// Consumes: Private or Public, FuncDecl or VarDecl
	// Produces: Member
	SEM_MEMBER_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		eat(stack, 2,
			FINAL_MEMBER,
			[]Kind{FINAL_PRIVATE, FINAL_PUBLIC},
			[]Kind{FINAL_FUNC_DECL, FINAL_VAR_DECL})
	},

	// 'impl' definition
	//
	// Consumes: Id, FuncDefList
	// Produces: ImplDef
	SEM_IMPL_DEF_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		eat(stack, 2, FINAL_IMPL_DEF, FINAL_ID, FINAL_FUNC_DEF_LIST)
	},

	// Function definition list for ImplDef
	//
	// Consumes: (FuncDefList, FuncDef) or (FuncDef) or ()
	// Produces: FuncDefList
	SEM_FUNCDEFLIST_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		list(stack, FINAL_FUNC_DEF_LIST, FINAL_FUNC_DEF)
	},

	// Consumes: Variable, AssignOp, ArithExpr or RelExpr
	// Produces: Assign
	SEM_ASSIGN_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		pullUp(stack, 1,
			FINAL_VARIABLE,
			FINAL_ASSIGN,
			[]Kind{FINAL_ARITH_EXPR, FINAL_REL_EXPR})
	},

	// Pushes an empty StatBlock on the stack. Use this to avoid listMerging
	// consecutive StatBlocks
	SEM_STATBLOCK_FRESH: func(stack *[]*ASTNode, tok Token) {
		pushNode(stack, FINAL_STATBLOCK, Token{})
	},

	// Consumes: RelExpr, StatBlock, StatBlock
	// Produces: If
	SEM_IF_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		eat(stack, 3, FINAL_IF, FINAL_REL_EXPR, FINAL_STATBLOCK, FINAL_STATBLOCK)
	},

	// Consumes: RelExpr, StatBlock
	// Produces: While
	SEM_WHILE_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		eat(stack, 2, FINAL_WHILE, FINAL_REL_EXPR, FINAL_STATBLOCK)
	},

	// Consumes: (StatBlock, some statement) or (some statement) or ()
	// Produces: StatBlock
	SEM_STATBLOCK_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		list(stack, FINAL_STATBLOCK, statementTypes...)
	},

	// Consumes: (FuncCallParamList, FuncCallParam) or (FuncCallParam) or ()
	// Produces: FuncCallParamList
	SEM_FUNC_CALL_PARAMLIST_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		list(stack, FINAL_FUNC_CALL_PARAMLIST, FINAL_FUNC_CALL_PARAM)
	},

	// Consumes: ArithExpr, RelExpr
	// Produces: FuncCallParam
	SEM_FUNC_CALL_PARAM_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		swapTop(stack, FINAL_FUNC_CALL_PARAM, FINAL_ARITH_EXPR, FINAL_REL_EXPR)
	},

	// Consumes: Subject, Id, ParamList,
	// Produces: FuncCall
	SEM_FUNC_CALL_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		eat(stack, 3, FINAL_FUNC_CALL, FINAL_SUBJECT, FINAL_ID, FINAL_FUNC_CALL_PARAMLIST)
	},

	// Consumes: Variable
	// Produces: Read
	SEM_READ_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		wrapTop(stack, FINAL_READ, FINAL_VARIABLE)
	},

	// Make a Subject node - the target of a Variable of FuncCall
	//
	// Consumes: Variable or FuncCall
	// Produces: Subject
	SEM_SUBJECT_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		lengthOrPanic(stack, 1)
		top := top(stack)
		switch top.Type {

		// Subjects are variables or function calls
		case FINAL_VARIABLE, FINAL_FUNC_CALL:
			wrapTop(stack, FINAL_SUBJECT, FINAL_VARIABLE, FINAL_FUNC_CALL)

		// If there is no valid subject, we will push an empty placeholder subject
		default:
			pushNode(stack, FINAL_SUBJECT, Token{})
		}
	},

	// Consumes: Subject, Id, IndexList
	// Produces: Variable
	SEM_VARIABLE_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		eat(stack, 3, FINAL_VARIABLE, FINAL_SUBJECT, FINAL_ID, FINAL_INDEXLIST)
	},

	SEM_INDEX_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		swapTop(stack, FINAL_INDEX, FINAL_ARITH_EXPR)
	},

	SEM_INDEXLIST_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		list(stack, FINAL_INDEXLIST, FINAL_INDEX)
	},

	SEM_DIMLIST_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		list(stack, FINAL_DIMLIST, FINAL_DIM)
	},

	SEM_DIM_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		swapTop(stack, FINAL_DIM, FINAL_INTNUM)
	},

	SEM_REL_EXPR_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		wrapTop(stack, FINAL_REL_EXPR, relOpTypes...)
	},

	SEM_REL_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		pullUp(stack, 1, FINAL_ARITH_EXPR, relOpTypes, FINAL_ARITH_EXPR)
	},

	SEM_MULTOP_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		pullUp(stack, 1,
			exprOperandTypes,
			[]Kind{FINAL_MULT, FINAL_DIV, FINAL_AND},
			exprOperandTypes)
	},

	SEM_ADDOP_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		pullUp(stack, 1,
			exprOperandTypes,
			[]Kind{FINAL_PLUS, FINAL_MINUS, FINAL_OR},
			exprOperandTypes)
	},

	SEM_EXPR_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		// ? noop for now
		// Just do a small type check
		lengthOrPanic(stack, 1)
		typeCheck(top(stack), FINAL_ARITH_EXPR, FINAL_REL_EXPR)
	},

	SEM_ARITH_EXPR_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		wrapTop(stack,
			FINAL_ARITH_EXPR,
			append(
				[]Kind{FINAL_TERM, FINAL_PLUS, FINAL_MINUS, FINAL_OR},
				exprOperandTypes...)...)
	},

	SEM_INTNUM_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		pushNode(stack, FINAL_INTNUM, tok)
	},

	SEM_FLOATNUM_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		pushNode(stack, FINAL_FLOATNUM, tok)
	},

	SEM_TERM_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		// ? noop for now

		// Just do a type check
		lengthOrPanic(stack, 1)
		typeCheck(top(stack), exprOperandTypes...)

		// Otherwise, wrap
		// wrapTop(
		// 	FINAL_TERM, stack,
		// 	append([]Kind{FINAL_FACTOR}, exprOperands...)...)
	},

	SEM_NEGATIVE_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		pushNode(stack, FINAL_NEGATIVE, tok)
	},

	SEM_POSITIVE_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		pushNode(stack, FINAL_POSITIVE, tok)
	},

	SEM_NOT_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		pushNode(stack, FINAL_NOT, tok)
	},

	SEM_FACTOR_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		lengthOrPanic(stack, 1)
		switch top := top(stack); top.Type {
		case FINAL_INTNUM,
			FINAL_FLOATNUM,
			FINAL_ARITH_EXPR,
			FINAL_VARIABLE,
			FINAL_FUNC_CALL:
			wrapTop(stack,
				FINAL_FACTOR,
				FINAL_INTNUM,
				FINAL_FLOATNUM,
				FINAL_ARITH_EXPR,
				FINAL_VARIABLE,
				FINAL_FUNC_CALL)

		case FINAL_FACTOR:
			// Then we need to consume two from the top
			lengthOrPanic(stack, 2)
			transformWithChildren(stack, 2, &ASTNode{Type: FINAL_FACTOR})

		default:
			typeCheck(top, FINAL_INTNUM) // Otherwise panic
		}
	},

	SEM_STATEMENT_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		// ? noop for now, just do a typeCheck
		lengthOrPanic(stack, 1)
		typeCheck(top(stack), statementTypes...)
	},

	SEM_FUNC_BODY_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		list(stack,
			FINAL_FUNC_BODY,
			append([]Kind{FINAL_VAR_DECL}, statementTypes...)...)
	},

	SEM_INTEGER_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		pushNode(stack, FINAL_INTEGER, tok)
	},

	SEM_FUNC_DEF_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		eat(stack, 4, FINAL_FUNC_DEF,
			FINAL_ID, FINAL_FUNC_DEF_PARAMLIST,
			FINAL_RETURNTYPE, FINAL_FUNC_BODY)
	},

	SEM_ID_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		pushNode(stack, FINAL_ID, tok)
	},

	SEM_REPT_PROG0_MAKESIBLING: func(stack *[]*ASTNode, tok Token) {
		list(
			stack,
			FINAL_STRUCT_OR_IMPL_OR_FUNC_LIST,
			FINAL_FUNC_DEF, FINAL_IMPL_DEF, FINAL_STRUCT_DECL)
	},

	SEM_WRITE_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		wrapTop(stack,
			FINAL_WRITE,
			FINAL_EXPR, FINAL_ARITH_EXPR, FINAL_REL_EXPR)
	},

	SEM_RETURN_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		wrapTop(stack,
			FINAL_RETURN,
			FINAL_EXPR, FINAL_ARITH_EXPR, FINAL_REL_EXPR)
	},

	// When this action is called, stack should look like:
	// [FINAL_STRUCT_OR_IMPL_OR_FUNC_LIST]
	SEM_PROG_MAKE_NODE: func(stack *[]*ASTNode, tok Token) {
		l := len(*stack)
		if l != 1 {
			panic(fmt.Errorf(
				"expected stack to be [FINAL_STRUCT_OR_IMPL_OR_FUNC_LIST] but got %v",
				*stack))
		}

		top := (*stack)[l-1]
		if top.Type != FINAL_STRUCT_OR_IMPL_OR_FUNC_LIST {
			panic(fmt.Errorf(
				"expected stack to be [FINAL_STRUCT_OR_IMPL_OR_FUNC_LIST] but got %v",
				*stack))
		}

		node := &ASTNode{
			Type:     FINAL_PROG,
			Children: []*ASTNode{top},
		}
		(*stack)[l-1] = node
	},

	SEM_FLOAT_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		pushNode(stack, FINAL_FLOAT, tok)
	},

	SEM_VOID_MAKENODE: func(stack *[]*ASTNode, tok Token) {
		pushNode(stack, FINAL_VOID, tok)
	},

	SEM_TYPE_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		wrapTop(stack,
			FINAL_TYPE,
			FINAL_INTEGER, FINAL_FLOAT, FINAL_ID)
	},

	SEM_RETURNTYPE_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		lengthOrPanic(stack, 1)
		switch top := top(stack); top.Type {
		case FINAL_TYPE:
			swapTop(stack, FINAL_RETURNTYPE, FINAL_TYPE)
		case FINAL_VOID:
			wrapTop(stack, FINAL_RETURNTYPE, FINAL_VOID)
		default:
			typeCheck(top, FINAL_TYPE, FINAL_VOID) // This will panic
		}
	},

	SEM_FPARAM_LIST_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		list(stack, FINAL_FUNC_DEF_PARAMLIST, FINAL_FUNC_DEF_PARAM)
	},

	SEM_FPARAM_MAKEFAMILY: func(stack *[]*ASTNode, tok Token) {
		eat(stack, 3, FINAL_FUNC_DEF_PARAM, FINAL_ID, FINAL_TYPE, FINAL_DIMLIST)
	},

	SEM_PLUS_MAKENODE:     func(s *[]*ASTNode, t Token) { pushNode(s, FINAL_PLUS, t) },
	SEM_MINUS_MAKENODE:    func(s *[]*ASTNode, t Token) { pushNode(s, FINAL_MINUS, t) },
	SEM_OR_MAKENODE:       func(s *[]*ASTNode, t Token) { pushNode(s, FINAL_OR, t) },
	SEM_MULT_MAKENODE:     func(s *[]*ASTNode, t Token) { pushNode(s, FINAL_MULT, t) },
	SEM_DIV_MAKENODE:      func(s *[]*ASTNode, t Token) { pushNode(s, FINAL_DIV, t) },
	SEM_AND_MAKENODE:      func(s *[]*ASTNode, t Token) { pushNode(s, FINAL_AND, t) },
	SEM_EQ_MAKENODE:       func(s *[]*ASTNode, t Token) { pushNode(s, FINAL_EQ, t) },
	SEM_NEQ_MAKENODE:      func(s *[]*ASTNode, t Token) { pushNode(s, FINAL_NEQ, t) },
	SEM_LT_MAKENODE:       func(s *[]*ASTNode, t Token) { pushNode(s, FINAL_LT, t) },
	SEM_GT_MAKENODE:       func(s *[]*ASTNode, t Token) { pushNode(s, FINAL_GT, t) },
	SEM_LEQ_MAKENODE:      func(s *[]*ASTNode, t Token) { pushNode(s, FINAL_LEQ, t) },
	SEM_GEQ_MAKENODE:      func(s *[]*ASTNode, t Token) { pushNode(s, FINAL_GEQ, t) },
	SEM_ASSIGNOP_MAKENODE: func(s *[]*ASTNode, t Token) { pushNode(s, FINAL_ASSIGN, t) },
	SEM_PUBLIC_MAKENODE:   func(s *[]*ASTNode, t Token) { pushNode(s, FINAL_PUBLIC, t) },
	SEM_PRIVATE_MAKENODE:  func(s *[]*ASTNode, t Token) { pushNode(s, FINAL_PRIVATE, t) },
}
