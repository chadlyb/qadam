package shared

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ISO CP 852
// We use unicode 2400 "symbol for NUL" for NUL (0), so it is printable
// We use unicode 2423 "Open Box" for NBSP (255), so it is printable
// We use unicode 00A5 "Yen" aka the paragraph sign for the section symbol (u00A7), so there aren't two.
// We use unicode 00AF "Macron" for the soft hyphen, so there aren't two.
const ctrlCharacters = "\u2400\u263a\u263b\u2665\u2666\u2663\u2660\u2022\u25D8\u25CB\u25D9\u2642\u2640\u266A\u266B\u263C\u25BA\u25C4\u2195\u203C\u00B6\u00A5\u25AC\u21A8\u2191\u2193\u2192\u2190\u221F\u2194\u25B2\u25BC"
const CharsetString = ctrlCharacters + " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~⌂ÇüéâäůćçłëŐőîŹÄĆÉĹĺôöĽľŚśÖÜŤťŁ×čáíóúĄąŽžĘę¬źČş«»░▒▓│┤ÁÂĚŞ╣║╗╝Żż┐└┴┬├─┼Ăă╚╔╩╦╠═╬¤đĐĎËďŇÍÎě┘┌█▄ŢŮ▀ÓßÔŃńňŠšŔÚŕŰýÝţ´\u00AF˝˛ˇ˘§÷¸°¨˙űŘř■\u2423"

var CharsetRunes = []rune(CharsetString)

var CharsetMapToByte = map[rune]byte{}

var czechHistogram = [256]float32{}
var czechHistogramLog = [256]float32{}
var uniformHistogram = [256]float32{}
var uniformHistogramLog = [256]float32{}

func init() {
	for k, v := range CharsetRunes {
		_, has := CharsetMapToByte[v]
		if has {
			fmt.Printf("%c redundant (%x vs %x)\n", v, k, CharsetMapToByte[v])

		}
		if !has {
			CharsetMapToByte[v] = byte(k)
		}
	}

	CharsetMapToByte['\r'] = '\r'
	CharsetMapToByte['\n'] = '\n'
	CharsetMapToByte['\t'] = '\t'

	// Unicode madness
	CharsetMapToByte['–'] = CharsetMapToByte['-']
	CharsetMapToByte['—'] = CharsetMapToByte['-']
	CharsetMapToByte['‘'] = CharsetMapToByte['\'']
	CharsetMapToByte['’'] = CharsetMapToByte['\'']

	czechBytes, err := FromString(czechSample)
	if err != nil {
		panic(fmt.Errorf("couldn't get czech histogram: %w", err))
	}

	const epsilon = 1e-8
	czechHistogram = NormalizedHistogramFromBytes(czechBytes)
	for i := 0; i < 256; i++ {
		czechHistogramLog[i] = float32(math.Log(float64(czechHistogram[i])))
		if czechHistogramLog[i] < epsilon {
			czechHistogramLog[i] = epsilon
		}
	}

	uniformHistogram = [256]float32{1.0 / 256}
	uniformHistogramLog = [256]float32{float32(math.Log(float64(1.0 / 256)))}
}

func UnescapeString(s string) (string, error) {
	var out strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			i++
			if i >= len(s) {
				return "", errors.New("trailing backslash in string")
			}
			switch s[i] {
			case 'n':
				out.WriteByte('\n')
			case 't':
				out.WriteByte('\t')
			case 'r':
				out.WriteByte('\r')
			case '\\':
				out.WriteByte('\\')
			case '"':
				out.WriteByte('"')
			case 'x':
				if i+2 >= len(s) {
					return "", errors.New("invalid \\x escape in string")
				}
				b, err := strconv.ParseUint(s[i+1:i+3], 16, 8)
				if err != nil {
					return "", fmt.Errorf("invalid hex in \\x escape: %v", err)
				}
				out.WriteByte(byte(b))
				i += 2
			default:
				return "", fmt.Errorf("unknown escape sequence: \\%c", s[i])
			}
		} else {
			out.WriteByte(s[i])
		}
	}
	return out.String(), nil
}

func ToString(bytes []byte) string {
	var b strings.Builder
	for _, v := range bytes {
		switch v {
		case '"':
			b.WriteString("\\\"")
		case '\n':
			b.WriteString("\\n")
		case '\t':
			b.WriteString("\\t")
		case '\\':
			b.WriteString("\\\\")
		default:
			b.WriteRune(CharsetRunes[v])
		}
	}
	return b.String()
}

func FromString(str string) ([]byte, error) {
	unescaped, err := UnescapeString(str)
	if err != nil {
		return nil, err
	}

	result := make([]byte, 0, len(unescaped))
	for _, v := range unescaped {
		n, ok := CharsetMapToByte[v]
		if !ok {
			return nil, fmt.Errorf("unrecognized character '%c' (%d)", v, int(v))
		}
		result = append(result, n)
	}

	return result, nil
}

// https://www.gutenberg.org/cache/epub/34225/pg34225.txt
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

func HistogramFromBytes(bytes []byte) [256]int {
	histogram := [256]int{}
	for _, b := range bytes {
		histogram[b]++
	}
	return histogram
}

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

func HistogramDotProduct(a [256]float32, b [256]float32) float32 {
	total := float32(0)
	for i := 0; i < 256; i++ {
		total += a[i] * b[i]
	}
	return total
}

func LogLikelihood(bytes []byte, refHistLog [256]float32) float64 {
	logLikelihood := 0.0
	for _, b := range bytes {
		logLikelihood += float64(refHistLog[b])
	}
	return logLikelihood
}

func IsLikelyCzechString(bytes []byte) bool {
	const minLength = 2
	if len(bytes) < minLength {
		return false
	}

	llCzech := LogLikelihood(bytes, czechHistogramLog)
	llUniform := LogLikelihood(bytes, uniformHistogramLog)
	normalizedDiff := (llCzech - llUniform) / float64(len(bytes))

	// Debug print the string and the normalized difference
	//fmt.Printf("'%v' normalized diff: %f\n", string(bytes), normalizedDiff)

	const threshold = 3.0
	return normalizedDiff > threshold
}

/*
func HeuristicIsHumanString(bytes []byte) bool {
	if len(bytes) < 2 {
		return false
	}
	spaceCount := 0
	ctrlCount := 0
	for _, v := range bytes {
		if v == ' ' {
			spaceCount++
		}
		if v < 32 {
			ctrlCount++
		}
	}
	if len(bytes) < 4 && ctrlCount > 0 {
		return false
	}

	if len(bytes) > 16 && spaceCount == 0 {
		return false
	}

	return true
}
*/
