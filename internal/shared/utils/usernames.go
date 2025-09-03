package utils

import (
	"math/rand"
	"strconv"
)

var adjectives = []string{
	"swift", "clever", "bright", "bold", "quick", "smart", "cool", "sharp",
	"brave", "calm", "happy", "lucky", "mighty", "noble", "royal", "wise",
	"epic", "super", "ultra", "mega", "blazing", "cosmic", "stellar", "quantum",
	"shadow", "golden", "silver", "crimson", "azure", "mystic", "electric", "digital",
	"fierce", "gentle", "silent", "loud", "wild", "tame", "free", "bound",
	"dark", "light", "fast", "slow", "hot", "cold", "smooth", "rough",
	"sleek", "rugged", "modern", "ancient", "future", "retro", "neon", "chrome",
	"velvet", "steel", "glass", "marble", "wooden", "plastic", "metal", "crystal",
	"frozen", "burning", "flying", "diving", "racing", "dancing", "singing", "roaring",
	"whispering", "screaming", "laughing", "crying", "smiling", "frowning", "glowing", "sparking",
	"spinning", "jumping", "crawling", "walking", "running", "hiding", "seeking", "finding",
	"broken", "fixed", "bent", "straight", "curved", "twisted", "simple", "complex",
	"tiny", "huge", "mini", "giant", "small", "large", "narrow", "wide",
	"thick", "thin", "deep", "shallow", "high", "low", "tall", "short",
	"young", "old", "new", "vintage", "fresh", "stale", "pure", "mixed",
}

var nouns = []string{
	"wolf", "eagle", "lion", "tiger", "dragon", "phoenix", "hawk", "bear",
	"storm", "thunder", "lightning", "comet", "star", "moon", "sun", "galaxy",
	"knight", "warrior", "hunter", "ranger", "wizard", "ninja", "samurai", "viking",
	"robot", "cyborg", "hacker", "coder", "pixel", "byte", "matrix", "circuit",
	"mountain", "ocean", "forest", "river", "flame", "crystal", "diamond", "steel",
	"cat", "dog", "fox", "rabbit", "deer", "horse", "dolphin", "shark",
	"bird", "fish", "snake", "turtle", "frog", "butterfly", "bee", "spider",
	"flower", "tree", "grass", "leaf", "branch", "root", "seed", "fruit",
	"rock", "stone", "sand", "clay", "mud", "dirt", "dust", "ash",
	"wind", "rain", "snow", "ice", "mist", "fog", "cloud", "rainbow",
	"sword", "shield", "arrow", "bow", "spear", "hammer", "axe", "blade",
	"crown", "ring", "gem", "jewel", "treasure", "coin", "gold", "silver",
	"book", "pen", "paper", "ink", "scroll", "letter", "word", "story",
	"music", "song", "dance", "art", "paint", "brush", "canvas", "sculpture",
	"car", "bike", "plane", "ship", "train", "rocket", "boat", "truck",
	"house", "tower", "castle", "bridge", "gate", "door", "window", "roof",
	"key", "lock", "chain", "rope", "thread", "wire", "cable", "net",
	"ball", "cube", "sphere", "triangle", "square", "circle", "oval", "diamond",
	"game", "toy", "puzzle", "riddle", "mystery", "secret", "code", "cipher",
	"dream", "wish", "hope", "fear", "joy", "love", "peace", "chaos",
	"elephant", "giraffe", "zebra", "hippo", "rhino", "cheetah", "leopard", "panther",
	"gorilla", "monkey", "chimp", "orangutan", "panda", "koala", "kangaroo", "sloth",
	"penguin", "flamingo", "parrot", "owl", "crow", "raven", "robin", "sparrow",
	"whale", "octopus", "squid", "jellyfish", "starfish", "seahorse", "crab", "lobster",
	"ant", "beetle", "cricket", "grasshopper", "dragonfly", "mosquito", "fly", "wasp",
	"rose", "lily", "tulip", "daisy", "sunflower", "orchid", "violet", "jasmine",
	"oak", "pine", "maple", "willow", "birch", "cedar", "palm", "bamboo",
	"apple", "banana", "orange", "grape", "cherry", "strawberry", "peach", "pear",
	"carrot", "potato", "tomato", "pepper", "onion", "garlic", "spinach", "broccoli",
	"bread", "cheese", "milk", "butter", "honey", "sugar", "salt", "pepper",
	"coffee", "tea", "juice", "water", "soda", "wine", "beer", "whiskey",
	"chair", "table", "bed", "sofa", "desk", "lamp", "mirror", "clock",
	"phone", "computer", "laptop", "tablet", "camera", "television", "radio", "speaker",
	"hammer", "screwdriver", "wrench", "drill", "saw", "nail", "screw", "bolt",
	"shirt", "pants", "dress", "skirt", "jacket", "coat", "hat", "shoes",
	"watch", "necklace", "bracelet", "earring", "glasses", "wallet", "purse", "bag",
	"planet", "asteroid", "meteor", "nebula", "universe", "cosmos", "void", "black hole",
	"mountain", "valley", "canyon", "cliff", "cave", "desert", "jungle", "swamp",
	"lake", "pond", "stream", "waterfall", "beach", "island", "volcano", "glacier",
	"castle", "palace", "mansion", "cottage", "cabin", "tent", "igloo", "hut",
	"temple", "church", "mosque", "shrine", "monument", "statue", "pillar", "arch",
	"garden", "park", "playground", "stadium", "arena", "theater", "museum", "library",
	"hospital", "school", "university", "factory", "office", "store", "market", "mall",
	"restaurant", "cafe", "bar", "hotel", "motel", "inn", "lodge", "resort",
	"train", "subway", "bus", "taxi", "helicopter", "jet", "spacecraft", "satellite",
	"sword", "dagger", "katana", "saber", "rapier", "scimitar", "claymore", "machete",
	"shield", "armor", "helmet", "gauntlet", "boots", "cloak", "cape", "robe",
	"wand", "staff", "orb", "amulet", "talisman", "charm", "pendant", "medallion",
	"potion", "elixir", "antidote", "poison", "medicine", "herb", "flower", "mushroom",
	"scroll", "tome", "grimoire", "spell", "curse", "blessing", "prayer", "chant",
	"crystal", "prism", "lens", "mirror", "window", "portal", "gateway", "passage",
	"labyrinth", "maze", "dungeon", "chamber", "vault", "treasury", "armory", "library",
	"forge", "anvil", "furnace", "cauldron", "crucible", "altar", "shrine", "temple",
}

func GenerateRandomUsername() string {
	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	num := rand.Intn(999) + 1

	return adj + "_" + noun + "_" + strconv.Itoa(num)
}
