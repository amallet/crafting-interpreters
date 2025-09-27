package main

type Stmt interface {
	Accept(visitor StmtVisitor) error 
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt *ExpressionStmt) error 
	VisitPrintStmt(stmt *PrintStmt) error 
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

type VarStmt struct {
	name Token
	initializer Expr
}

func (v *VarStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitVarStmt(v)
}



