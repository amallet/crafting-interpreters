package main

type Stmt interface {
	Accept(visitor StmtVisitor) error
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt *ExpressionStmt) error
	VisitFunctionStmt(stmt *FunctionStmt) error
	VisitClassStmt(stmt *ClassStmt) error 
	VisitIfStmt(stmt *IfStmt) error
	VisitPrintStmt(stmt *PrintStmt) error
	VisitWhileStmt(stmt *WhileStmt) error
	VisitReturnStmt(stmt *ReturnStmt) error
	VisitBlockStmt(stmt *BlockStmt) error
	VisitVarStmt(stmt *VarStmt) error
}

type ExpressionStmt struct {
	expression Expr
}

func (e *ExpressionStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitExpressionStmt(e)
}

type ClassStmt struct {
	className Token
	methods []*FunctionStmt
}

func (c *ClassStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitClassStmt(c)
}


type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

type FunctionStmt struct {
	functionName Token
	isGetter bool 
	params       []Token
	body         []Stmt
}

func (f *FunctionStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitFunctionStmt(f)
}

func (i *IfStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitIfStmt(i)
}

type PrintStmt struct {
	expression Expr
}

func (s *PrintStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitPrintStmt(s)
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func (w *WhileStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitWhileStmt(w)
}

type ReturnStmt struct {
	keyword     Token
	returnValue Expr
}

func (r *ReturnStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitReturnStmt(r)
}

type BlockStmt struct {
	statements []Stmt
}

func (b *BlockStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitBlockStmt(b)
}

type VarStmt struct {
	variable    Token
	initializer Expr
}

func (v *VarStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitVarStmt(v)
}
