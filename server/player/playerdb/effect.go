package playerdb

import "github.com/df-mc/dragonfly/server/entity/effect"

func effectsToData(effects []effect.Effect) []jsonEffect {
	data := make([]jsonEffect, len(effects))
	for key, eff := range effects {
		id, ok := effect.ID(eff.Type())
		if !ok {
			continue
		}
		data[key] = jsonEffect{
			ID:              id,
			Duration:        eff.Duration(),
			Level:           eff.Level(),
			Ambient:         eff.Ambient(),
			ParticlesHidden: eff.ParticlesHidden(),
			Infinite:        eff.Infinite(),
		}
	}
	return data
}

func dataToEffects(data []jsonEffect) []effect.Effect {
	effects := make([]effect.Effect, 0, len(data))
	for _, d := range data {
		e, ok := effect.ByID(d.ID)
		if !ok {
			continue
		}
		switch eff := e.(type) {
		case effect.LastingType:
			var lasting effect.Effect
			switch {
			case d.Ambient:
				lasting = effect.NewAmbient(eff, d.Level, d.Duration)
			case d.Infinite:
				lasting = effect.NewInfinite(eff, d.Level)
			default:
				lasting = effect.New(eff, d.Level, d.Duration)
			}

			if d.ParticlesHidden {
				lasting = lasting.WithoutParticles()
			}
			effects = append(effects, lasting)
		default:
			effects = append(effects, effect.NewInstant(eff, d.Level))
		}
	}
	return effects
}
