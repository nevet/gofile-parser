package component

import (
	"fmt"
)

type GoFile struct {
	Package   string
	Import    []string
	Const     []string
	Var       *Definition
	Type      *Definition
	Functions []*Function
}

func (goFile *GoFile) Dump() {
	fmt.Println("File Info:")

	fmt.Println("\nPacakge Name:\n" + goFile.Package)

	fmt.Println("\nImport:")
	if len(goFile.Import) == 0 {
		fmt.Println("Not Found")
	} else {
		for _, v := range goFile.Import {
			fmt.Println(v)
		}
	}

	fmt.Println("\nConst:")
	if len(goFile.Const) == 0 {
		fmt.Println("Not Found")
	} else {
		for _, v := range goFile.Const {
			fmt.Println(v)
		}
	}

	fmt.Println("\nVar:")
	if len(*goFile.Var) == 0 {
		fmt.Println("Not Found")
	} else {
		goFile.Var.Dump()
	}

	fmt.Println("\nType:")
	if len(*goFile.Type) == 0 {
		fmt.Println("Not Found")
	} else {
		goFile.Type.Dump()
	}

	fmt.Println("\nFunc:")
	if len(goFile.Functions) == 0 {
		fmt.Println("Not Found")
	} else {
		for _, v := range goFile.Functions {
			v.Dump()
		}
	}
}
