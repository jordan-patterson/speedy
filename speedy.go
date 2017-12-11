package main

import(
  "fmt"
  "os"
  "os/user"
  "github.com/jordan-patterson/promptui"
  "encoding/json"
  "io/ioutil"
  "os/exec"
)

type Directory struct{
  Alias string `json:"alias"`
  Path string `json:"path"`
}
type Directories struct{
  Directories []Directory `json:"directories"`
}

func exists(path string)bool{
  //reports if the given path exists or not
  if _,err:=os.Stat(path);err==nil{
    //exists
    return true
  }
  //not exists
  return false
}
func getFilePath()string{
  //returns the filepath where the saved data are stored
  usr,err:=user.Current()
  if err!=nil{
    fmt.Println(err)
  }
  filepath := usr.HomeDir+"/bin/speedy/dirs.json"
  return filepath
}

func getAllDirs()(map[string]string,[]string){
  //returns a dictionary/map of aliases to paths
  var dirs Directories
  jsonFile,err := os.Open(getFilePath())
  if err!=nil{
    fmt.Println(err)
  }
  //close file
  defer jsonFile.Close()
  bytes,_:= ioutil.ReadAll(jsonFile)
  err=json.Unmarshal(bytes,&dirs)
  if err!=nil{
    fmt.Println(err)
  }
  size:=len(dirs.Directories)
  //create string array of length size
  aliases:=make([]string,size)
  directories:=make(map[string]string)
  for i:=0;i<size;i++{
    alias:=dirs.Directories[i].Alias
    directories[alias]=dirs.Directories[i].Path
    aliases[i]=alias
  }
  return directories,aliases
}

func promptDirs(){
  //simple arrow prompt where user can choose a destination directory or cancel
  _,aliases:=getAllDirs()
  opts:=append(aliases,"ADD","REMOVE","CANCEL")
  prompt:=promptui.Select{
    Label:"Select An Alias",
    Items:opts,
  }
  _,result,err:=prompt.Run()
  if err!=nil{
    fmt.Println(err)
    return
  }
  if result!="CANCEL" && result!="ADD" && result!="REMOVE"{
    //path:=dirs[result]
    changeDir(result)
  }else if result=="ADD"{
    getNewDirectory()
  }else if result=="REMOVE"{
    removeDir()
  }
}

func getNewDirectory(){
  confirm:=promptui.Select{
    Label:"Are you sure?",
    Items:[]string{"YES","NO","CANCEL"},
  }
  var alias,path string
  for{
    fmt.Print("Enter an alias or CTRL + C to cancel: ")
    fmt.Scanln(&alias)
    fmt.Print("Enter the path or CTRL + C to cancel: ")
    for{
      fmt.Scanln(&path)
      //check if path exists
      if(exists(path)){
        break
      }
      fmt.Print("Enter a path that EXISTS or CTRL + C to cancel: ")
    }
    fmt.Printf("%s : %s\n",alias,path)
    _,res,err:=confirm.Run()
    if err!=nil{
      fmt.Println(err)
      return
    }
    if(res=="YES"){
      break
    }else if(res=="CANCEL"){
      return
    }
  }
  saveDir(alias,path)
}

func saveDir(alias,path string){
  //adds new pair to loaded map from getAllDirs, then updates file
  dirs,aliases:=getAllDirs()
  dirs[alias]=path
  //update alias array
  newAliases:=append(aliases,alias)
  updateDirs(dirs,newAliases)
}

func updateDirs(directories map[string]string,aliases []string){
  //overwrites .json file on each call with passed map
  size:=len(aliases)
  //construct Directory array to store in main struct
  dirs:=make([]Directory,size)
  for i:=0;i<size;i++{
    dirs[i].Alias=aliases[i]
    dirs[i].Path=directories[aliases[i]]
  }
  //now make Directories struct container
  jsonDirs:=Directories{dirs}
  //fmt.Println(jsonDirs)
  jsonData,err := json.Marshal(jsonDirs)
  if err!=nil{
    fmt.Println(err)
    return
  }
  //open file for writing
  file,err2:=os.Create(getFilePath())
  if err2!=nil{
    fmt.Println(err2)
    return
  }
  defer file.Close()
  //fmt.Println(string(jsonData))
  //write json data to file
  file.Write(jsonData)
  file.Close()
}

func defined(name string)bool{
  //reports whether or not alias(name) is defined
  dirs,_:=getAllDirs()
  if _,ok :=dirs[name]; ok{
    return true
  }
  return false
}

func changeDir(alias string){
  //go to directory mapped to passed alias
  dirs,_:=getAllDirs()
  path:=dirs[alias]
  fmt.Printf("(: Opening : %s\n",path)
  os.Chdir(path)
  //msg:=":) Teleported to : "+path
  cmd:=exec.Command("gnome-terminal")
  err:=cmd.Run()
  if err!=nil{
    fmt.Println(err)
  }
  msg:="Opened to "+path+" from here."
  fmt.Println(msg)
}

/**************MAIN*******************/
func main(){
//initialize directory where data will be stored
  if !exists(getFilePath()){
    os.Mkdir(getFilePath(),0700)
  }
  argLength:=len(os.Args)
  if(len(os.Args)>1){
    alias:=os.Args[argLength-1]
    if(defined(alias)){
      changeDir(alias)
    }else{
      fmt.Println("That alias is not defined. :(")
    }
  }else{
    promptDirs()
  }
}

/************OTHER FUNCTIONS********/
func removeDir(){
  //removes pair from loaded dictionary/map based on alias
  dirs,aliases:=getAllDirs()
  opts:=append(aliases,"CANCEL")
  prompt:=promptui.Select{
    Label:"Select An Alias To Remove",
    Items:opts,
  }
  _,res,err:=prompt.Run()
  if err!=nil{
    fmt.Println(err)
    return
  }
  if res!="CANCEL"{
    delete(dirs,res)
    //update alias array
    newAliases:=make([]string,len(aliases)-1)
    i:=0
    for k,_ := range dirs{
      newAliases[i]=k
      i++
    }
    updateDirs(dirs,newAliases)
  }
  return
}
