package types

type Cardinality int

type CardinalityRange struct {
	Min Cardinality
	Max Cardinality
}

const (
	CardinalityZero Cardinality = 0
	CardinalityOne  Cardinality = 1
	CardinalityMany Cardinality = 2
)

func CardinalityNone() CardinalityRange {
	return CardinalityRange{
		Min: CardinalityZero,
		Max: CardinalityZero,
	}
}

func CardinalityZeroToOne() CardinalityRange {
	return CardinalityRange{
		Min: CardinalityZero,
		Max: CardinalityOne,
	}
}

func CardinalityZeroToMany() CardinalityRange {
	return CardinalityRange{
		Min: CardinalityZero,
		Max: CardinalityMany,
	}
}

func CardinalityOneToOne() CardinalityRange {
	return CardinalityRange{
		Min: CardinalityOne,
		Max: CardinalityOne,
	}
}

func CardinalityOneToMany() CardinalityRange {
	return CardinalityRange{
		Min: CardinalityOne,
		Max: CardinalityMany,
	}
}

func CardinalityManyToMany() CardinalityRange {
	return CardinalityRange{
		Min: CardinalityMany,
		Max: CardinalityMany,
	}
}
