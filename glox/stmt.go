package main

type Stmt interface {
	Accept(visitor StmtVisitor) error 
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt *ExpressionStmt) error 
	VisitPrintStmt(stmt *PrintStmt) error 
	VisitBlockStmt(stmt *BlockStmt) error 
	VisitVarStmt(stmt *VarStmt) error 
}

type ExpressionStmt struct {
	expression Expr 
}

func (e* ExpressionStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitExpressionStmt(e)
}

type PrintStmt struct {
	expression Expr
}

func (s *PrintStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitPrintStmt(s)
}

type BlockStmt struct {
	statements []Stmt
}

func (b *BlockStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitBlockStmt(b)
}

type VarStmt struct {
	name Token
	initializer Expr
}

func (v *VarStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitVarStmt(v)
}



