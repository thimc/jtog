# jtog

*J*SON *to* *G*o is a command line tool that converts JSON to valid Go
source code.

## Usage

	Usage: jtog [ -i=bool ] [ -l=bool ] [ -o=bool ] [ file ... ]
	If no file path(s) are specified as flags then data from standard input is assumed.
	
	 -i	indent using spaces
	 -l	inline type defintions (default true)
	 -o	appends "omitempty" to the json tag

## License

MIT
