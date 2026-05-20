package models

const (
	MCPBlockCount         = 8
	MCPDimensionsPerBlock = 8
	MCPVectorSize         = MCPBlockCount * MCPDimensionsPerBlock
)

var RoleWeights = map[string][]float64{
	"engineer": {
		0.20,
		0.30,
		0.10,
		0.10,
		0.20,
		0.05,
		0.03,
		0.02,
	},
	"manager": {
		0.10,
		0.05,
		0.15,
		0.15,
		0.03,
		0.25,
		0.02,
		0.25,
	},
	"startup": {
		0.25,
		0.15,
		0.02,
		0.10,
		0.25,
		0.05,
		0.15,
		0.03,
	},
}

var BlockNames = []string{
	"Execution Drive",
	"Cognitive Processing",
	"Professional Values",
	"Emotional Resilience",
	"Learning Velocity",
	"Collaboration Depth",
	"Operational Rhythm",
	"Team Positioning",
}

var BlockDifferences = map[string]string{
	"Execution Drive":      "Different levels of initiative, ownership, and focus on achieving results",
	"Cognitive Processing": "Different ways of thinking, analyzing information, and making decisions",
	"Professional Values":  "Possible differences in core professional values and priorities",
	"Emotional Resilience": "Different ways of handling stress, pressure, and emotional responses",
	"Learning Velocity":    "Differences in speed of learning, adaptability, and applying new knowledge",
	"Collaboration Depth":  "Different approaches to teamwork, communication, and shared problem-solving",
	"Operational Rhythm":   "Potential mismatch in work pace, structure, and execution style",
	"Team Positioning":     "Different preferences in roles and contribution within team dynamics",
}
