package operatorstate

func Initialize() *F2SOperatorState {
	return &F2SOperatorState{
		IsMaster:       false,
		KnownOperators: make([]F2SKnownOperator, 0),
	}
}
