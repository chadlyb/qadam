package shared

import (
	"fmt"
	"math"
	"sort"
)

// LanguageModel represents a simple language model based on byte frequencies.
// Hist contains normalized frequencies (sum = 1). HistLog stores ln(prob) for fast log-likelihood.
type LanguageModel struct {
	Hist    [256]float32
	HistLog [256]float32
}

// newLanguageModel builds a model from raw byte counts (already normalized to frequencies).
func newLanguageModel(hist [256]float32) LanguageModel {
	const epsilon = 1e-8
	var m LanguageModel
	m.Hist = hist
	for i := 0; i < 256; i++ {
		p := hist[i]
		if p < epsilon {
			p = epsilon // avoid log(0)
		}
		m.HistLog[i] = float32(math.Log(float64(p)))
	}
	return m
}

// NewLanguageModelFromBytes builds a model from a sample byte slice.
func NewLanguageModelFromBytes(data []byte) LanguageModel {
	return newLanguageModel(NormalizedHistogramFromBytes(data))
}

// NewUniformLanguageModel returns a model where every byte is equally likely.
func NewUniformLanguageModel() LanguageModel {
	var hist [256]float32
	for i := 0; i < 256; i++ {
		hist[i] = 1.0 / 256.0
	}
	return newLanguageModel(hist)
}

// LogLikelihood computes the log-likelihood of data under this model.
func (m LanguageModel) LogLikelihood(data []byte) float64 {
	ll := 0.0
	for _, b := range data {
		ll += float64(m.HistLog[b])
	}
	return ll
}

// TrimmedLogLikelihood computes the log-likelihood using only the best-matching characters.
// It drops the worst-matching trimPercent of characters to make the detection more robust.
func (m LanguageModel) TrimmedLogLikelihood(data []byte, trimPercent float64) float64 {
	if len(data) == 0 {
		return 0.0
	}

	// Calculate individual character log-likelihoods
	charScores := make([]float64, len(data))
	for i, b := range data {
		charScores[i] = float64(m.HistLog[b])
	}

	// Sort scores in descending order (best matches first)
	sort.Sort(sort.Reverse(sort.Float64Slice(charScores)))

	// Calculate how many characters to trim
	trimCount := int(float64(len(data)) * trimPercent)
	if trimCount >= len(data) {
		trimCount = len(data) - 1 // Keep at least one character
	}

	// Sum only the best-matching characters
	total := 0.0
	for i := 0; i < len(data)-trimCount; i++ {
		total += charScores[i]
	}

	return total
}

// JSDivergence returns the Jensen-Shannon divergence between two models (symmetric, always finite).
// https://en.wikipedia.org/wiki/Jensen%E2%80%93Shannon_divergence
func (m LanguageModel) JSDivergence(other LanguageModel) float64 {
	// Convert to float64 for accuracy.
	const epsilon = 1e-8
	var kl1, kl2 float64
	for i := 0; i < 256; i++ {
		p := float64(m.Hist[i])
		q := float64(other.Hist[i])
		mVal := 0.5 * (p + q)
		if p > 0 {
			kl1 += p * math.Log((p+epsilon)/(mVal+epsilon))
		}
		if q > 0 {
			kl2 += q * math.Log((q+epsilon)/(mVal+epsilon))
		}
	}
	return 0.5*kl1 + 0.5*kl2
}

// HistogramFromBytes creates a raw histogram from byte data.
func HistogramFromBytes(bytes []byte) [256]int {
	histogram := [256]int{}
	for _, b := range bytes {
		histogram[b]++
	}
	return histogram
}

// NormalizedHistogramFromBytes creates a normalized histogram from byte data.
func NormalizedHistogramFromBytes(bytes []byte) [256]float32 {
	histogram := HistogramFromBytes(bytes)
	total := 0
	for _, v := range histogram {
		total += v
	}
	normalized := [256]float32{}
	for i := 0; i < 256; i++ {
		normalized[i] = float32(histogram[i]) / float32(total)
	}
	return normalized
}

// Czech text sample for histogram generation
const czechSample = `
MRTVÝ DŮM.

Trestnice naše stála na kraji pevnosti u samého pevnostního
náspu. Pohlédne-li časem člověk skrze štěrbinu ve plotě na boží
svět, nespatři-li tam něco nového, uvidí všeho všudy proužek
oblohy, pak vysoký, ze země nasypaný val, porostlý buřením, a po
valu sem tam ve dne v noci se procházejí stráže. A tu napadne
člověku, že minou celá léta a opět se přijde podívat štěrbinou ve
plotě a opět uvidí tentýž val, tytéž stráže a tentýž uzounký
proužek nebe — ne toho nebe, jež se klene nad trestnicí, nýbrž
jiného, dalekého, volného nebe. Představte si veliký dvůr, v
podobě nepravidelného šestiúhelníku, dvě stě kroků dlouhý, půl
druhého sta kroků široký, celý kolkolem ohrazený vysokým týnem,
to jest plotem z vysokých břeven, stojmo zakopaných hluboko do
země, těsně k sobě přiléhajících svými boky, upevněných příčními
plaňkami a nahoře přiostřených: toť vnější ohrada trestnice. V
jedné stěně ohrady jsou upravena pevná vrata vždy uzavřená,
vždycky — ve dne i v noci — hlídaná strážemi. Otvírala se jen na
daný rozkaz, aby propustila trestance, jdoucí na práci. Za těmi
vraty byl světlý, svobodný svět; tam žili lidé jako všude. Ale
zde, uvnitř ohrady měli o onom světě představy takové, jako o
nějaké nesplnitelné pohádce. Zde byl svůj zvláštní svět, jenž
neměl nikde nic sobě podobného; zde byly své zvláštní zákony, své
obleky, své mravy a obyčeje, a za živa mrtvý dům; život, jako
nikde jinde, a též lidé zvláštní. A právě tento zvláštní kout
hodlám popsati.

Vstoupíte-li do ohrady, spatříte v ní několik budov. Po obou
stranách širokého vnitřního dvora táhnou se dvě dlouhá, přízemní
stavení, sroubená ze dřeva. To jsou kasárny. V nich bydlí
trestanci, rozdělení dle trestů. Dále v pozadí ohrady ještě jedno
takové stavení: je to kuchyně, rozdělená pro dvě skupiny. Dále
ještě jedno stavení, kde jsou pod společnou střechou umístěny
sklepy, sýpky a kůlny. Prostředek dvora je prázdný a tvoří rovné,
dosti rozsáhlé prostranství. Zde se trestanci staví do řad,
provádí se kontrola a vyvolávají se jména ráno, v poledne a
večer, někdy však ještě mimo to několikrát za den, podle toho,
jak jsou stráže podezřívavé a jak umějí rychle počítati. Kolkolem
mezi budovami a plotem zůstává ještě dosti veliké prostranství.
Zde v zadu za budovami procházívají se rádi za doby odpočinku
někteří z vězňů, kteří nemilují společnosti a jsou povahy
zasmušilejší; chráněni před cizími zraky obírají se svými
myšlénkami. Setkal-li jsem se s některými z nich za takové
procházky, pozoroval jsem se zálibou jich zasmušilé, znamenané
tváře, a hleděl jsem uhodnouti, o čem asi přemýšlejí. Byl mezi
nimi jeden deportovaný, jehož zamilovaným zaměstnáním za svobodné
chvíle bylo počítati břevna. Bylo jich na půl druhého tisíce, a
všechna je měl spočítána i poznamenána. Každé břevno mu znamenalo
jeden den; každodenně vynechával jedno břevno a takovým způsobem
dle množství zbývajících břeven mohl názorně poznati, kolik dní
mu ještě zbývá zůstávati v trestnici, než vyprší lhůta jeho
nucené práce. Míval vždy nelíčenou radost, když se přiblížil ke
konci některé strany šestiúhelníka. Zbývalo mu čekati ještě mnoho
let; ale v trestnici měl dosti.času, aby se naučil trpělivosti.

Viděl jsem kdysi, jak se loučil se soudruhy jeden trestanec, jenž
strávil v káznici dvacet let a konečně odcházel na svobodu. Byli
lidé, kteří se pamatovali, jak poprvé vstoupil do trestnice
mladý, bezstarostný, nedbající ani svého zločinu, ani svého
trestu. A odcházel šedým starcem, s tváří zasmušilou a
truchlivou. Mlčky obešel všech šest našich kasáren. Když vstoupil
do kasárny, poznamenal se křížem směrem k svatým obrazům, a pak
se hluboko, až po pás poklonil soudruhům, žádaje na nich, aby
zapomněli, učinil-li jim co zlého.

Také si vzpomínám, jak kdysi jednoho vězně, někdy zámožného
sibiřského sedláka, zavolali jednou pod večer ke vratům. Před půl
rokem obdržel zprávu, že bývalá jeho žena se vdala a těžce se
proto zarmoutil. A nyní sama ona zajela ke trestnici, dala ho
zavolati a podala mu almužnu. Pohovořili spolu dvě minuty, oba si
poplakali a rozloučili se na věky. Viděl jsem jeho tvář, když se
vracel do kasárny... Ano, v tom místě bylo možná naučiti se
trpělivosti.

Když se smrákalo, všechny nás odváděli do kasáren, kde nás
zavírali na celou noc. Bylo mně vždy těžko vraceti se ze dvora do
naší kasárny. Byla to dlouhá, nízká, dusná síň, mdle osvětlená
lojovými svícemi, s obtížným, dusivým zápachem. Nechápu nyní, jak
jsem v ní strávil deset let. Na narách*) měl jsem svá tři prkna:
to bylo celé moje místo. Na těchto narách bylo jen v naší
světnici umístěno třicet lidí. V zimě zavíralo se záhy. Trvávalo
čtyry hodiny, než všichni usnuli. Před tím hluk, křik, chechtot,
nadávky, řinkot řetězů, dým, smrad, oholené hlavy, cejchované
tváře, hadrovité šaty, všechno osypané nadávkami, zdarebačené...
ano, tuhý má člověk život! Člověk jest bytost ke všemu
přivykající a myslím, to právě že jest nejlepší jeho definice.

Bylo nás v trestnici celkem dvě stě padesát osob — počet ten se
skoro neměnil. Jedni přicházeli, druhým vyšla lhůta i odcházeli,
třetí umírali. A co tu bylo za lidi! Myslím, že tu měla své
zástupce každá gubernie, každá oblasť celého Ruska. Byli mezi
trestanci jinorodci, bylo i několik horalů z Kavkazu, Vše to se
dělilo podle velikosti zločinů a následovně i podle počtu let,
přisouzených za zločin. Možno se domýšleti, že není na světě
zločinu, aby zde nebyl mel svého zástupce. Hlavní element všeho
obyvatelstva káznice tvořili deportovaní trestanci stavu
občanského, po`

// English text sample for histogram generation
const englishSample = `
THE HOUSE OF THE DEAD.

Our prison stood on the edge of the fortress, right by the rampart. If a person sometimes looked through a crack in the fence at the world outside, unless something new had appeared, all he would see was a strip of sky, then a high embankment covered with weeds, and on the embankment, day and night, guards would walk back and forth. And it would occur to a person that years would pass, and he would come again to look through the crack in the fence, and again he would see the same embankment, the same guards, and the same narrow strip of sky—not the sky that arches over the prison, but another, distant, free sky...`

// Global language models
var czechModel LanguageModel
var uniformModel LanguageModel
var EnglishModel LanguageModel

// initLanguageModels initializes the global language models
func initLanguageModels() {
	czechBytes, err := FromString(czechSample)
	if err != nil {
		panic(fmt.Errorf("couldn't get czech histogram: %w", err))
	}

	czechModel = NewLanguageModelFromBytes(czechBytes)

	englishBytes, err := FromString(englishSample)
	if err != nil {
		panic(fmt.Errorf("couldn't get english histogram: %w", err))
	}
	EnglishModel = NewLanguageModelFromBytes(englishBytes)

	uniformModel = NewUniformLanguageModel()
}

// IsLikelyCzechString determines if text is likely Czech using language model comparison
func IsLikelyCzechString(bytes []byte) bool {
	const minLength = 3
	if len(bytes) < minLength {
		return false
	}

	var normalizedDiff float64

	if len(bytes) < 10 {
		// For short strings, use non-trimmed scoring with lower threshold
		llCzech := czechModel.LogLikelihood(bytes)
		llUniform := uniformModel.LogLikelihood(bytes)
		normalizedDiff = (llCzech - llUniform) / float64(len(bytes))
	} else {
		// For longer strings, use trimmed log-likelihood
		const trimPercent = 0.15
		llCzech := czechModel.TrimmedLogLikelihood(bytes, trimPercent)
		llUniform := uniformModel.TrimmedLogLikelihood(bytes, trimPercent)

		// Calculate effective length (after trimming)
		effectiveLength := len(bytes) - int(float64(len(bytes))*trimPercent)
		if effectiveLength <= 0 {
			effectiveLength = 1
		}

		normalizedDiff = (llCzech - llUniform) / float64(effectiveLength)
	}

	// Use different thresholds based on string length
	var threshold float64
	if len(bytes) < 6 {
		threshold = -1.5 // Very lenient threshold for very short strings
	} else if len(bytes) < 10 {
		threshold = -0.5 // Lenient threshold for short strings
	} else {
		threshold = 0.8 // Higher threshold for longer strings
	}

	// Add debug output for strings that are close to the threshold
	if normalizedDiff > threshold-0.5 && normalizedDiff < threshold+0.5 {
		if len(bytes) < 10 {
			fmt.Printf("DEBUG: Language model score for \"%s\": %.3f (threshold: %.1f, short string)\n",
				ToString(bytes), normalizedDiff, threshold)
		} else {
			fmt.Printf("DEBUG: Language model score for \"%s\": %.3f (threshold: %.1f, trimmed %.1f%%)\n",
				ToString(bytes), normalizedDiff, threshold, 0.15*100)
		}
	}

	return normalizedDiff > threshold
}

// IsLikelyEnglishString determines if text is likely English using language model comparison
func IsLikelyEnglishString(bytes []byte) bool {
	const minLength = 3
	if len(bytes) < minLength {
		return false
	}

	var normalizedDiff float64

	if len(bytes) < 10 {
		// For short strings, use non-trimmed scoring with lower threshold
		llEnglish := EnglishModel.LogLikelihood(bytes)
		llUniform := uniformModel.LogLikelihood(bytes)
		normalizedDiff = (llEnglish - llUniform) / float64(len(bytes))
	} else {
		// For longer strings, use trimmed log-likelihood
		const trimPercent = 0.15
		llEnglish := EnglishModel.TrimmedLogLikelihood(bytes, trimPercent)
		llUniform := uniformModel.TrimmedLogLikelihood(bytes, trimPercent)

		// Calculate effective length (after trimming)
		effectiveLength := len(bytes) - int(float64(len(bytes))*trimPercent)
		if effectiveLength <= 0 {
			effectiveLength = 1
		}

		normalizedDiff = (llEnglish - llUniform) / float64(effectiveLength)
	}

	// Use different thresholds based on string length
	var threshold float64
	if len(bytes) < 6 {
		threshold = -1.5 // Very lenient threshold for very short strings
	} else if len(bytes) < 10 {
		threshold = -0.5 // Lenient threshold for short strings
	} else {
		threshold = 0.8 // Higher threshold for longer strings
	}

	if normalizedDiff > threshold-0.5 && normalizedDiff < threshold+0.5 {
		if len(bytes) < 10 {
			fmt.Printf("DEBUG: English model score for \"%s\": %.3f (threshold: %.1f, short string)\n",
				ToString(bytes), normalizedDiff, threshold)
		} else {
			fmt.Printf("DEBUG: English model score for \"%s\": %.3f (threshold: %.1f, trimmed %.1f%%)\n",
				ToString(bytes), normalizedDiff, threshold, 0.15*100)
		}
	}
	return normalizedDiff > threshold
}

// IsLikelyHumanLanguage returns true if the string is likely Czech or English
func IsLikelyHumanLanguage(bytes []byte) bool {
	return IsLikelyCzechString(bytes) || IsLikelyEnglishString(bytes)
}

// DebugTrimmedComparison compares trimmed vs non-trimmed scores for debugging
func DebugTrimmedComparison(bytes []byte) {
	if len(bytes) < 5 {
		return
	}

	// Non-trimmed scores
	llCzech := czechModel.LogLikelihood(bytes)
	llUniform := uniformModel.LogLikelihood(bytes)
	normalizedDiff := (llCzech - llUniform) / float64(len(bytes))

	// Trimmed scores (15%)
	const trimPercent = 0.15
	llCzechTrimmed := czechModel.TrimmedLogLikelihood(bytes, trimPercent)
	llUniformTrimmed := uniformModel.TrimmedLogLikelihood(bytes, trimPercent)
	effectiveLength := len(bytes) - int(float64(len(bytes))*trimPercent)
	if effectiveLength <= 0 {
		effectiveLength = 1
	}
	normalizedDiffTrimmed := (llCzechTrimmed - llUniformTrimmed) / float64(effectiveLength)

	fmt.Printf("DEBUG COMPARISON for \"%s\":\n", ToString(bytes))
	fmt.Printf("  Non-trimmed: %.3f\n", normalizedDiff)
	fmt.Printf("  Trimmed (15%%): %.3f\n", normalizedDiffTrimmed)
	fmt.Printf("  Difference: %.3f\n", normalizedDiffTrimmed-normalizedDiff)
}
