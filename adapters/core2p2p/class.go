package core2p2p

import (
	"fmt"

	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/p2p/gen"
	"github.com/NethermindEth/juno/utils"
)

func AdaptClass(class core.Class) *gen.Class {
	if class == nil {
		return nil
	}

	hash, err := class.Hash()
	if err != nil {
		panic(fmt.Errorf("failed to hash %t: %w", class, err))
	}

	switch v := class.(type) {
	case *core.Cairo0Class:
		return &gen.Class{
			Class: &gen.Class_Cairo0{
				Cairo0: &gen.Cairo0Class{
					Abi:          string(v.Abi),
					Externals:    utils.Map(v.Externals, adaptEntryPoint),
					L1Handlers:   utils.Map(v.L1Handlers, adaptEntryPoint),
					Constructors: utils.Map(v.Constructors, adaptEntryPoint),
					Program:      v.Program,
				},
			},
			Domain:    0, // todo(kirill) recheck
			ClassHash: AdaptHash(hash),
		}
	case *core.Cairo1Class:
		return &gen.Class{
			Class: &gen.Class_Cairo1{
				Cairo1: &gen.Cairo1Class{
					Abi: v.Abi,
					EntryPoints: &gen.Cairo1EntryPoints{
						Externals:    utils.Map(v.EntryPoints.External, adaptSierra),
						L1Handlers:   utils.Map(v.EntryPoints.L1Handler, adaptSierra),
						Constructors: utils.Map(v.EntryPoints.Constructor, adaptSierra),
					},
					Program:              utils.Map(v.Program, AdaptFelt),
					ContractClassVersion: v.SemanticVersion,
				},
			},
			Domain:    0, // todo(kirill) recheck
			ClassHash: AdaptHash(hash),
		}
	default:
		panic(fmt.Errorf("unsupported cairo class %T (version=%d)", v, class.Version()))
	}
}

func adaptSierra(e core.SierraEntryPoint) *gen.SierraEntryPoint {
	return &gen.SierraEntryPoint{
		Index:    e.Index,
		Selector: AdaptFelt(e.Selector),
	}
}

func adaptEntryPoint(e core.EntryPoint) *gen.EntryPoint {
	return &gen.EntryPoint{
		Selector: AdaptFelt(e.Selector),
		Offset:   e.Offset.Uint64(),
	}
}
