# jtog

*J*SON *to* *G*o is a command line tool that converts JSON to valid Go
source code.

## Usage

	Usage: jtog [ -l=bool ] [ -o=bool ] [ file ... ]
	If no file path(s) are specified as flags then data from standard input is assumed.
	
	 -l	inline type defintions (default true)
	 -o	appends "omitempty" to the json tag

### Examples

	jtog file1.json file2.json ...

or

	jtog <file1.json

## License

MIT
