package pkg

const (
	//terminals
	digit = "[0-9]"

	zero           = "quot"
	quarter        = "#188"
	half           = "#189"
	halfAndQuarter = "#190"

	//non-terminals
	numberRegex = "[1-9]?" + digit + "+"
	mesureRegex = numberRegex + " ?\""

	floatingCommas = "(" + zero + "|" + quarter + "|" + half + "|" + halfAndQuarter + ")"

	quote              = "<ut type=\"unknown\" x=\"" + digit + "+" + "\">&amp;amp;" + floatingCommas + ";</ut>"
	imperialMatchRegex = numberRegex + " ?" + quote + "(" + quote + ")?"
	sourceExp          = "<source.*>.*" + quote + ".*?" + "</source>"

	sourceWithoutQuote = "<source>.*" + numberRegex + ".*</source>"
)
