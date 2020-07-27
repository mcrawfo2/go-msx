package skel

const (
	archetypeKeyApp         = "app"
	archetypeKeyBeat        = "beat"
	archetypeKeyServicePack = "sp"
)

type Archetype struct {
	Key         string
	DisplayName string
	Generators  []string
}

type Archetypes []Archetype

func (a Archetypes) Key(index int) string {
	return a[index].Key
}

func (a Archetypes) DisplayNames() []string {
	var result []string
	for _, archetype := range a {
		result = append(result, archetype.DisplayName)
	}
	return result
}

func (a Archetypes) Generators(key string) []string {
	for _, arch := range a {
		if arch.Key == key {
			return arch.Generators
		}
	}
	return nil
}

var archetypes = Archetypes{
	{
		Key:         archetypeKeyApp,
		DisplayName: "Generic Microservice",
		Generators: []string{
			"generate-migrate",
		},
	},
	{
		Key:         archetypeKeyBeat,
		DisplayName: "Beat",
		Generators: []string{
			"generate-domain-beats",
		},
	},
	{
		Key:         archetypeKeyServicePack,
		DisplayName: "Service Pack Microservice",
		Generators: []string{
			"generate-migrate",
			"generate-service-pack",
		},
	},
}
