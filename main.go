package main
	import (
		"fmt"
		"strings"
		"io/ioutil"
		"encoding/json"
		"container/list"
		"reflect"

		"log"
	)


type BDFileSymbol struct{
	key string
	value interface{}
}


type BDFileElement struct {
	handler string
	value string
	listOfSymbol list.List
}


 func (BDFE *BDFileElement) getValue() string{
	 return BDFE.value
 }

 func (BDFE *BDFileElement) setValue ( value string ){
		BDFE.value = value
	 	var tmp_map interface{}
	 	err := json.Unmarshal([]byte(value), &tmp_map)
	 if err != nil{
		 fmt.Println("Ошибка разбиения на символы")
		 //TODO-ME Добавить обработку некорректного JSON Объекта
	 }

	 keys := reflect.ValueOf(tmp_map).MapKeys()
	 //TODO-ME Добавить Определение типов для послудующей конвертации значений
	 for _,v := range keys{
		 var tmpSymbol BDFileSymbol
		 tmpSymbol.key =v.Interface().(string)
		 tmpSymbol.value = tmp_map.(map[string]interface{})[tmpSymbol.key]
		 BDFE.listOfSymbol.PushBack(tmpSymbol)

	 }



 }

func (BDFE *BDFileElement) setHandler( handler string ){
	BDFE.handler = handler
}
//BLOCK

func newBlock() *BDFileBlock  {
	p:= new(BDFileBlock)
	p.next = nil
	for k := range p.elements{
		p.elements[k].handler= "0"

	}
	//fmt.Println("CREATE BLOCK")
	return p
}


type BDFileBlock  struct {
	next *BDFileBlock
	index int
	elements [4]BDFileElement
}
func (BDFB *BDFileBlock) getValue() string{
	 BlockValue := ""

	for k := range BDFB.elements{
		if k == BDFB.index{
			break
		}
		if len(BlockValue) > 0 {
			BlockValue+=","
		}

		 BlockValue += BDFB.elements[k].getValue()
	 }

	return BlockValue
}
func (BDFB *BDFileBlock) setValue(value string) (int){
	writeIndex,full := BDFB.isFull()
	if full {

		return -1
	}else{

		BDFB.elements[writeIndex].setValue(value)

		BDFB.index++
		return 1
	}

}

func (BDFB *BDFileBlock) isFull() (int, bool){
	index := BDFB.index
	if index == 4 {
		return index, true
	}

	return index, false
}

func newFile() *BDFile{
	p := new(BDFile)
	blck := newBlock()

	p.first =blck
	p.last = p.first
	p.length = 1
	return p
}

type BDFileIndex struct {
	eIndex int
	bIndex int
}
func (BDFI *BDFileIndex) getIndex(index int) *BDFileIndex{
	BDFI.eIndex = index%4
	BDFI.bIndex = ((index-BDFI.eIndex)/4)
	return BDFI

}

func (BDF *BDFile)getElementByIndex (i int) string{
	index := new(BDFileIndex).getIndex(i)
	p := BDF.first
	for i:= 0; i < index.bIndex; i++  {
		if p == nil{
			fmt.Println("OUT OF RANGE")
			log.Fatalln("FUCK FUCK FUCK")
		}
		p=p.next
	}
	return p.elements[index.eIndex].value
}

type BDFile struct {
	first *BDFileBlock
	last *BDFileBlock
	length int
	fileindex BDFileIndex
}

func (BDF *BDFile) getValue() string{
	p := BDF.first
	result := ""
	for p != nil{
		if len(result)>0 {
			result+=","

		}
		result += p.getValue()
		p = p.next
	}
	return "["+result+"]"
}
func (BDF *BDFile) setValue(val string) {
	p := BDF.last
	err := p.setValue( val )
	if ( -1 == err ){

		blck := newBlock()
		blck.setValue( val )


		BDF.last.next = blck
		BDF.last = blck
		BDF.length++
	}
}

func (BDF *BDFile) setFromFile( value string ){
	value = value[1:len(value)-1]
	parseParam := strings.SplitAfter(value,"},")
	for _,v := range parseParam {
		v := strings.Replace(v,"},","}",-1)
		BDF.setValue(v)
	}

}
//TODO-ME Реализовать поиск в файле по ключу и значению ключа ( "КЛЮЧ":"ЗНАЧЕНИЕ" )

func main() {

	/*
	File := newFile()

	File.setValue("{Name: Anton ,Age : 23, city: Moscow}")
	File.setValue("{Name: Andrey ,Age : 34, city: Kiev}")
	File.setValue("{Name: Elza ,Age : 20, city: Penza}")
	File.setValue("{Name: Lara ,Age : 54, city: Stalinsk}")
	File.setValue("{Name: Anton1 ,Age : 22, city: Moscow}")
	File.setValue("{Name: Andrey1 ,Age : 33, city: Kiev}")
	File.setValue("{Name: Elza1 ,Age : 19, city: Penza}")
	File.setValue("{Name: Lara1 ,Age : 53, city: Stalinsk}")
	File.setValue("{Name: Voctory, Age:20, city: Novokuznetsk}")
	File.setValue("{Name: Lana, Age:20, city: Novokuznetsk}")
	File.setValue("{Name: OLGA, Age:20, city: Novokuznetsk}")
	File.setValue("{Name: VITYA, Age:20, city: Novokuznetsk}")
	File.setValue("{Name: EDIK, Age:20, city: Novokuznetsk}")

	*/

//	fmt.Println(File.getValue())

	FileTest := newFile()



	dat, err := ioutil.ReadFile("Nomenclature.txt")
	if err != nil{
		fmt.Println("Can't Read file!!!")
	}
	FileTest.setFromFile( string( dat ) )

	//fmt.Println(FileTest.getValue())
	//fmt.Println(File.length)
	fmt.Println(FileTest.length)
	fmt.Println(FileTest.getElementByIndex(5))

}
